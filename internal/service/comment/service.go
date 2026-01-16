package comment

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/repository/comment"
	"go-twitter/internal/repository/user"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID, postID int64, req dto.CreateCommentRequest) (int64, int, error)
	GetCommentByID(ctx context.Context, id int64) (*dto.CommentResponse, int, error)
	GetCommentsByPostID(ctx context.Context, postID int64, page, pageSize int) (*dto.CommentsResponse, int, error)
	UpdateComment(ctx context.Context, userID, commentID int64, req dto.UpdateCommentRequest) (int, error)
	DeleteComment(ctx context.Context, userID, commentID int64) (int, error)
}

type commentService struct {
	commentRepo comment.CommentRepository
	userRepo    user.UserRepository
}

func NewService(commentRepo comment.CommentRepository, userRepo user.UserRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}
