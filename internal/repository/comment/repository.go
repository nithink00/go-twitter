package comment

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.CommentModel) (int64, error)
	GetCommentByID(ctx context.Context, id int64) (*model.CommentModel, error)
	GetCommentsByPostID(ctx context.Context, postID int64, offset, limit int) ([]*model.CommentModel, int64, error)
	UpdateComment(ctx context.Context, comment *model.CommentModel) error
	DeleteComment(ctx context.Context, id int64) error
	GetCommentLikesCount(ctx context.Context, commentID int64) (int, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}
