package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetComment(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	comment, status, err := h.commentService.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
