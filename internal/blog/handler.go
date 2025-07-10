package blog

import (
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	blogs := r.Group("/blogs")
	{
		blogs.GET("", h.ListBlogs)
		blogs.GET("/:id", h.GetBlog)
		blogs.POST("", h.CreateBlog)
		blogs.PUT("/:id", h.UpdateBlog)
		blogs.DELETE("/:id", h.DeleteBlog)
		blogs.POST("/:id/like", h.LikeBlog)
		blogs.DELETE("/:id/like", h.UnlikeBlog)
		blogs.POST("/:id/bookmark", h.BookmarkBlog)
		blogs.DELETE("/:id/bookmark", h.UnbookmarkBlog)
	}
}

func (h *Handler) CreateBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Set author ID from authenticated user
	if userID, exists := c.Get("userID"); exists {
		blog.AuthorID = userID.(string)
	}

	if err := h.service.CreateBlog(c.Request.Context(), &blog); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, blog)
}

func (h *Handler) GetBlog(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	blog, err := h.service.GetBlogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (h *Handler) ListBlogs(c *gin.Context) {
	page := utils.GetQueryInt64(c, "page", 1)
	limit := utils.GetQueryInt64(c, "limit", 10)

	filter := bson.M{}
	if tag := c.Query("tag"); tag != "" {
		filter["tags"] = tag
	}
	if authorID := c.Query("author"); authorID != "" {
		filter["author_id"] = authorID
	}

	blogs, err := h.service.ListBlogs(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, blogs)
}

func (h *Handler) UpdateBlog(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	blog.ID = id
	if err := h.service.UpdateBlog(c.Request.Context(), &blog); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (h *Handler) DeleteBlog(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if err := h.service.DeleteBlog(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LikeBlog(c *gin.Context) {
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

	if err := h.service.LikeBlog(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UnlikeBlog(c *gin.Context) {
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

	if err := h.service.UnlikeBlog(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) BookmarkBlog(c *gin.Context) {
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

	if err := h.service.BookmarkBlog(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UnbookmarkBlog(c *gin.Context) {
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

	if err := h.service.UnbookmarkBlog(c.Request.Context(), id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}
