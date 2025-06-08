package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"github.com/dksensei/letsnormalizeit/internal/config"
	"github.com/dksensei/letsnormalizeit/internal/db"
	"github.com/dksensei/letsnormalizeit/internal/middleware"
	"github.com/dksensei/letsnormalizeit/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Firebase Auth service
	authService, err := auth.NewService(&cfg.Firebase)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}

	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close(context.Background())

	// Initialize Redis connection
	redis, err := db.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize repositories
	userRepo := user.NewRepository(mongodb)

	// Initialize services
	userService := user.NewService(userRepo, authService)

	// Initialize handlers
	authHandler := auth.NewHandler(authService, userService)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	// Cleanup old entries every 5 minutes
	rateLimiter.Cleanup(5 * time.Minute)

	// Setup Gin router
	router := gin.New() // Use New() instead of Default() to customize middleware

	// Use our custom recovery middleware
	router.Use(middleware.Recovery())

	// Use logger middleware
	router.Use(gin.Logger())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Server.AllowOrigins}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Apply rate limiting to all routes
	router.Use(rateLimiter.RateLimit())

	// Register auth routes
	authHandler.Register(router)

	// Public routes
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

	// Protected routes (require authentication)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.POST("/blogs", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Created blog by user: %s", uid)})
		})

		protected.POST("/blogs/:id/like", func(c *gin.Context) {
			id := c.Param("id")
			uid, _ := c.Get("uid")

			err := userService.ToggleLike(c.Request.Context(), uid.(string), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Blog %s like toggled by user %s", id, uid),
			})
		})

		protected.POST("/blogs/:id/bookmark", func(c *gin.Context) {
			id := c.Param("id")
			uid, _ := c.Get("uid")

			err := userService.ToggleBookmark(c.Request.Context(), uid.(string), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Blog %s bookmark toggled by user %s", id, uid),
			})
		})

		protected.POST("/comments", func(c *gin.Context) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Comment created by user: %s", uid)})
		})

		protected.GET("/user/bookmarks", func(c *gin.Context) {
			uid, _ := c.Get("uid")

			// Get user with bookmarks
			userData, err := userService.GetUserByID(c.Request.Context(), uid.(string))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user":      userData.Email,
				"bookmarks": userData.Bookmarks,
			})
		})

		protected.GET("/user/profile", func(c *gin.Context) {
			uid, _ := c.Get("uid")

			userData, err := userService.GetUserByID(c.Request.Context(), uid.(string))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":              userData.ID,
				"name":            userData.Name,
				"email":           userData.Email,
				"photo_url":       userData.PhotoURL,
				"created_at":      userData.CreatedAt,
				"bookmarks_count": len(userData.Bookmarks),
				"likes_count":     len(userData.Likes),
			})
		})

		protected.PUT("/user/profile", func(c *gin.Context) {
			uid, _ := c.Get("uid")

			var input struct {
				Name string `json:"name"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			userData, err := userService.UpdateUserProfile(c.Request.Context(), uid.(string), input.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":         userData.ID,
				"name":       userData.Name,
				"email":      userData.Email,
				"photo_url":  userData.PhotoURL,
				"updated_at": userData.UpdatedAt,
			})
		})
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(authService), middleware.AdminOnly(authService))
	{
		admin.GET("/users", func(c *gin.Context) {
			// In a real implementation, we would call userService.ListUsers()
			// with pagination parameters
			c.JSON(http.StatusOK, gin.H{"message": "List of users (admin only)"})
		})

		// Route to set admin privileges (would require additional service method)
		admin.POST("/users/:id/set-admin", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Set admin privileges for user: %s", id)})
		})
	}

	// Create a context that listens for signals to gracefully shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server running on port %s", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
