package user

import (
	"go-twitter/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, refreshToken, status, err := h.userService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
