package like

import (
	"go-twitter/internal/dto"
	"go-twitter/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LikeComment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	status, err := h.likeService.LikeComment(c.Request.Context(), int64(userID), commentID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusConflict {
		c.JSON(status, gin.H{"error": "comment already liked"})
		return
	}

	c.JSON(http.StatusCreated, dto.LikeResponse{Message: "comment liked successfully"})
}

func (h *Handler) UnlikeComment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	status, err := h.likeService.UnlikeComment(c.Request.Context(), int64(userID), commentID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusNotFound {
		c.JSON(status, gin.H{"error": "comment not liked yet"})
		return
	}

	c.JSON(http.StatusOK, dto.LikeResponse{Message: "comment unliked successfully"})
}

func (h *Handler) GetCommentLikesCount(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	count, status, err := h.likeService.GetCommentLikesCount(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LikesCountResponse{Count: count})
}
