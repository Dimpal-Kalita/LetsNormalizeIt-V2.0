package user

import (
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		utils.Error("Registration failed: %v", err)
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		utils.Error("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) GetProfile(c *gin.Context) {
	uid, _ := c.Get("uid")

	user, err := h.service.GetUserByID(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"photo_url":       user.PhotoURL,
		"created_at":      user.CreatedAt,
		"bookmarks_count": len(user.Bookmarks),
		"likes_count":     len(user.Likes),
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	uid, _ := c.Get("uid")

	var input struct {
		Name string `json:"name" binding:"required,min=2,max=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateUserProfile(c.Request.Context(), uid.(string), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"photo_url":  user.PhotoURL,
		"updated_at": user.UpdatedAt,
	})
}

func (h *Handler) GetBookmarks(c *gin.Context) {
	uid, _ := c.Get("uid")

	bookmarks, err := h.service.GetBookmarks(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bookmarks": bookmarks,
	})
}

func (h *Handler) ToggleLike(c *gin.Context) {
	uid, _ := c.Get("uid")
	blogID := c.Param("id")

	err := h.service.ToggleLike(c.Request.Context(), uid.(string), blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "like toggled"})
}

func (h *Handler) ToggleBookmark(c *gin.Context) {
	uid, _ := c.Get("uid")
	blogID := c.Param("id")

	err := h.service.ToggleBookmark(c.Request.Context(), uid.(string), blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bookmark toggled"})
}
