package post

import (
	"go-twitter/internal/middleware"
	"go-twitter/internal/service/post"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	api         *gin.Engine
	validate    *validator.Validate
	postService post.PostService
	authMiddleware *middleware.AuthMiddleware
}

func NewHandler(api *gin.Engine, validate *validator.Validate, postService post.PostService, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		api:         api,
		validate:    validate,
		postService: postService,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) RouteList() {
	postGroup := h.api.Group("/posts")
	{
		postGroup.GET("", h.GetPosts)
		postGroup.GET("/:post_id", h.GetPost)

		postGroup.Use(h.authMiddleware.RequireAuth())
		{
			postGroup.POST("", h.CreatePost)
			postGroup.PUT("/:post_id", h.UpdatePost)
			postGroup.DELETE("/:post_id", h.DeletePost)
		}
	}
}
