package comment

import (
	"go-twitter/internal/dto"
	"go-twitter/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateComment(c *gin.Context) {
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

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.commentService.UpdateComment(c.Request.Context(), int64(userID), commentID, req)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status == http.StatusForbidden {
		c.JSON(status, gin.H{"error": "you can only update your own comments"})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateCommentResponse{ID: commentID})
}
