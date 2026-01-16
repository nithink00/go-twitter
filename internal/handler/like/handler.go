package like

import (
	"go-twitter/internal/middleware"
	"go-twitter/internal/service/like"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	api            *gin.Engine
	likeService    like.LikeService
	authMiddleware *middleware.AuthMiddleware
}

func NewHandler(api *gin.Engine, likeService like.LikeService, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		api:            api,
		likeService:    likeService,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) RouteList() {
	postLikesGroup := h.api.Group("/posts/:post_id/likes")
	postLikesGroup.Use(h.authMiddleware.RequireAuth())
	{
		postLikesGroup.POST("", h.LikePost)
		postLikesGroup.DELETE("", h.UnlikePost)
		postLikesGroup.GET("/count", h.GetPostLikesCount)
	}

	commentLikesGroup := h.api.Group("/comments/:comment_id/likes")
	commentLikesGroup.Use(h.authMiddleware.RequireAuth())
	{
		commentLikesGroup.POST("", h.LikeComment)
		commentLikesGroup.DELETE("", h.UnlikeComment)
		commentLikesGroup.GET("/count", h.GetCommentLikesCount)
	}
}
