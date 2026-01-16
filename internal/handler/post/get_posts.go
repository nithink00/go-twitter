package post

import (
	"go-twitter/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	userIDStr := c.Query("user_id")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var posts *dto.PostsResponse
	var status int
	var err error

	if userIDStr != "" {
		userID, parseErr := strconv.ParseInt(userIDStr, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		posts, status, err = h.postService.GetPostsByUserID(c.Request.Context(), userID, page, pageSize)
	} else {
		posts, status, err = h.postService.GetPosts(c.Request.Context(), page, pageSize)
	}

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
