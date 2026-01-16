package post

import (
	"go-twitter/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeletePost(c *gin.Context) {
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

	status, err := h.postService.DeletePost(c.Request.Context(), int64(userID), postID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusForbidden {
		c.JSON(status, gin.H{"error": "you can only delete your own posts"})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}
