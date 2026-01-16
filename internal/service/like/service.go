package like

import (
	"context"
	"go-twitter/internal/repository/like"
)

type LikeService interface {
	LikePost(ctx context.Context, userID, postID int64) (int, error)
	UnlikePost(ctx context.Context, userID, postID int64) (int, error)
	GetPostLikesCount(ctx context.Context, postID int64) (int, int, error)

	LikeComment(ctx context.Context, userID, commentID int64) (int, error)
	UnlikeComment(ctx context.Context, userID, commentID int64) (int, error)
	GetCommentLikesCount(ctx context.Context, commentID int64) (int, int, error)
}

type likeService struct {
	likeRepo like.LikeRepository
}

func NewService(likeRepo like.LikeRepository) LikeService {
	return &likeService{
		likeRepo: likeRepo,
	}
}
