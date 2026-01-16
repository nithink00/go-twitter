package comment

import (
	"go-twitter/internal/middleware"
	"go-twitter/internal/service/comment"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	api            *gin.Engine
	validate       *validator.Validate
	commentService comment.CommentService
	authMiddleware *middleware.AuthMiddleware
}

func NewHandler(api *gin.Engine, validate *validator.Validate, commentService comment.CommentService, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		api:            api,
		validate:       validate,
		commentService: commentService,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) RouteList() {
	postCommentsGroup := h.api.Group("/posts/:post_id/comments")
	{
		postCommentsGroup.GET("", h.GetComments)

		postCommentsGroup.Use(h.authMiddleware.RequireAuth())
		{
			postCommentsGroup.POST("", h.CreateComment)
		}
	}

	commentsGroup := h.api.Group("/comments")
	{
		commentsGroup.GET("/:comment_id", h.GetComment)

		commentsGroup.Use(h.authMiddleware.RequireAuth())
		{
			commentsGroup.PUT("/:comment_id", h.UpdateComment)
			commentsGroup.DELETE("/:comment_id", h.DeleteComment)
		}
	}
}
