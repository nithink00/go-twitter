package post

import (
	"go-twitter/internal/dto"
	"go-twitter/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePost(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, status, err := h.postService.CreatePost(c.Request.Context(), int64(userID), req)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreatePostResponse{ID: postID})
}
