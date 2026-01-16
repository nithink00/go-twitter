package like

import (
	"go-twitter/internal/dto"
	"go-twitter/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LikePost(c *gin.Context) {
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

	status, err := h.likeService.LikePost(c.Request.Context(), int64(userID), postID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusConflict {
		c.JSON(status, gin.H{"error": "post already liked"})
		return
	}

	c.JSON(http.StatusCreated, dto.LikeResponse{Message: "post liked successfully"})
}

func (h *Handler) UnlikePost(c *gin.Context) {
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

	status, err := h.likeService.UnlikePost(c.Request.Context(), int64(userID), postID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusNotFound {
		c.JSON(status, gin.H{"error": "post not liked yet"})
		return
	}

	c.JSON(http.StatusOK, dto.LikeResponse{Message: "post unliked successfully"})
}

func (h *Handler) GetPostLikesCount(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	count, status, err := h.likeService.GetPostLikesCount(c.Request.Context(), postID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LikesCountResponse{Count: count})
}
