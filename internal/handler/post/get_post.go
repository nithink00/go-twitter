package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, status, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}
