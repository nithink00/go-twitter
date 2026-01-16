package user

import (
	"go-twitter/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	api *gin.Engine
	validate *validator.Validate
	userService user.UserService
}
func NewHandler(api *gin.Engine, validate *validator.Validate, userService user.UserService) *Handler {
	return &Handler{
		api:         api,
		validate:    validate,
		userService: userService,
	}
}

func (h *Handler) RouteList() {
	authGroup := h.api.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh", h.RefreshToken)
		authGroup.POST("/logout", h.Logout)
	}

	userGroup := h.api.Group("/users")
	{
		userGroup.GET("/:id", h.GetUser)
	}
}