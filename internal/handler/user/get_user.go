package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, status, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
