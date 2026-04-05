package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"github.com/dksensei/letsnormalizeit/internal/cache"
	"github.com/dksensei/letsnormalizeit/internal/config"
	"github.com/dksensei/letsnormalizeit/internal/db"
	"github.com/dksensei/letsnormalizeit/internal/features/user"
	"github.com/dksensei/letsnormalizeit/internal/middleware"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	loggerConfig := map[string]interface{}{
		"level":              cfg.Logger.Level,
		"encoding":           cfg.Logger.Encoding,
		"output_paths":       cfg.Logger.OutputPaths,
		"error_output_paths": cfg.Logger.ErrorOutputPaths,
	}
	utils.InitWithConfig(loggerConfig)
	defer utils.Sync()

	utils.Info("Starting LetsNormalizeIt-V2.0 server")
	utils.Info("Auth type: %s", cfg.Auth.Type)

	var mongoDB *db.MongoDB
	var err error

	mongoDB, err = db.NewMongoDB(&cfg.MongoDB)
	if err != nil {
		utils.Fatal("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close(context.Background())

	var cacheClient cache.Cache
	if cfg.Cache.Enabled {
		if cfg.Cache.Type == "redis" {
			redisCache, err := cache.NewRedisCache(&cfg.Redis, true)
			if err != nil {
				utils.Warn("Failed to connect to Redis, falling back to memory cache: %v", err)
				cacheClient = cache.NewMemoryCache()
			} else {
				cacheClient = redisCache
			}
		} else {
			cacheClient = cache.NewMemoryCache()
		}
		utils.Info("Cache enabled: %s", cfg.Cache.Type)
	} else {
		utils.Info("Cache disabled")
	}

	var jwtService *auth.Service
	if cfg.Auth.Type == "jwt" {
		jwtService = auth.NewService(
			cfg.Auth.JWT.Secret,
			cfg.Auth.JWT.AccessExpiry,
			cfg.Auth.JWT.RefreshExpiry,
		)
	}

	var authService auth.AuthService
	var userRepo user.Repository

	userRepo = user.NewRepository(mongoDB.Database)
	authService = auth.NewAuthService(jwtService, userRepo)

	userService := user.NewService(userRepo, authService)
	userHandler := user.NewHandler(userService)

	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	rateLimiter.Cleanup(5 * time.Minute)

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(gin.Logger())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Server.AllowOrigins}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	router.Use(rateLimiter.RateLimit())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "2.0",
			"auth":    cfg.Auth.Type,
			"cache":   cacheClient.IsEnabled(),
		})
	})

	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", userHandler.Register)
		authGroup.POST("/login", userHandler.Login)
		authGroup.POST("/refresh", userHandler.RefreshToken)
	}

	public := router.Group("/api/v1")
	{
		public.GET("/blogs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "List of blogs"})
		})

		public.GET("/blogs/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Blog with ID: %s", id)})
		})

		public.GET("/blogs/:id/comments", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Comments for blog ID: %s", id)})
		})
	}

	protected := router.Group("/api/v1")
	if jwtService != nil {
		jwtMiddleware := middleware.NewJWTMiddleware(jwtService)
		protected.Use(jwtMiddleware.Authenticate())

		protected.POST("/blogs", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Created blog by user: %s", uid)})
		})

		protected.POST("/blogs/:id/like", userHandler.ToggleLike)
		protected.POST("/blogs/:id/bookmark", userHandler.ToggleBookmark)

		protected.POST("/comments", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Comment created by user: %s", uid)})
		})

		protected.GET("/user/bookmarks", userHandler.GetBookmarks)
		protected.GET("/user/profile", userHandler.GetProfile)
		protected.PUT("/user/profile", userHandler.UpdateProfile)

		admin := protected.Group("/admin")
		admin.Use(jwtMiddleware.RequireAdmin())
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "List of users (admin only)"})
			})

			admin.POST("/users/:id/set-admin", func(c *gin.Context) {
				id := c.Param("id")
				c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Set admin privileges for user: %s", id)})
			})
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		utils.Info("Server running on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown: %v", err)
	}

	if cacheClient != nil {
		cacheClient.Close()
	}

	utils.Info("Server exited")
}
