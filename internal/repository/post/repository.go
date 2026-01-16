package post

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.PostModel) (int64, error)
	GetPostByID(ctx context.Context, id int64) (*model.PostModel, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*model.PostModel, error)
	GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.PostModel, error)
	GetPostsCount(ctx context.Context) (int64, error)
	UpdatePost(ctx context.Context, post *model.PostModel) error
	DeletePost(ctx context.Context, id int64) error
	GetPostWithUserInfo(ctx context.Context, id int64) (*model.PostModel, string, error)
	GetPostsWithUserInfo(ctx context.Context, limit, offset int) ([]*model.PostModel, []string, error)
}

type postRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}
