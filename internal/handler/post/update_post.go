package post

import (
	"go-twitter/internal/dto"
	"go-twitter/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdatePost(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.postService.UpdatePost(c.Request.Context(), int64(userID), postID, req)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusForbidden {
		c.JSON(status, gin.H{"error": "you can only update your own posts"})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UpdatePostResponse{ID: postID})
}
