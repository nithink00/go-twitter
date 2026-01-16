package like

import (
	"context"
	"database/sql"
)

type LikeRepository interface {
	LikePost(ctx context.Context, postID, userID int64) error
	UnlikePost(ctx context.Context, postID, userID int64) error
	IsPostLiked(ctx context.Context, postID, userID int64) (bool, error)
	GetPostLikesCount(ctx context.Context, postID int64) (int, error)

	LikeComment(ctx context.Context, commentID, userID int64) error
	UnlikeComment(ctx context.Context, commentID, userID int64) error
	IsCommentLiked(ctx context.Context, commentID, userID int64) (bool, error)
	GetCommentLikesCount(ctx context.Context, commentID int64) (int, error)
}

type likeRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) LikeRepository {
	return &likeRepository{
		db: db,
	}
}
