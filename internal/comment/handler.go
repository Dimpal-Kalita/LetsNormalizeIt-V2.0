package comment

import (
	"errors"
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	comments := r.Group("/comments")
	{
		comments.GET("/blog/:blogId", h.ListComments)
		comments.GET("/:id", h.GetComment)
		comments.POST("", h.CreateComment)
		comments.PUT("/:id", h.UpdateComment)
		comments.DELETE("/:id", h.DeleteComment)
		comments.POST("/:id/like", h.LikeComment)
		comments.DELETE("/:id/like", h.UnlikeComment)
	}
}

func (h *Handler) CreateComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Set user ID from authenticated user
	if userID, exists := c.Get("userID"); exists {
		comment.UserID = userID.(string)
	} else {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(errors.New("unauthorized")))
		return
	}

	if err := h.service.CreateComment(c.Request.Context(), &comment); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetComment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	comment, err := h.service.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if comment == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *Handler) ListComments(c *gin.Context) {
	blogID, err := primitive.ObjectIDFromHex(c.Param("blogId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	page := utils.GetQueryInt64(c, "page", 1)
	limit := utils.GetQueryInt64(c, "limit", 10)

	comments, err := h.service.ListCommentsByBlogID(c.Request.Context(), blogID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *Handler) UpdateComment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Verify ownership
	existingComment, err := h.service.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if existingComment == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		return
	}

	userID, exists := c.Get("userID")
	if !exists || existingComment.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	comment.ID = id
	comment.UserID = userID.(string)
	if err := h.service.UpdateComment(c.Request.Context(), &comment); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Verify ownership
	comment, err := h.service.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if comment == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		return
	}

	userID, exists := c.Get("userID")
	if !exists || comment.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	if err := h.service.DeleteComment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LikeComment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	if err := h.service.LikeComment(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UnlikeComment(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	if err := h.service.UnlikeComment(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}
