package post

import (
	"context"
	"database/sql"
	"go-twitter/internal/config"
	"go-twitter/internal/dto"
	"go-twitter/internal/repository/post"
)

type PostService interface {
	CreatePost(ctx context.Context, userID int64, req dto.CreatePostRequest) (int64, int, error)
	GetPostByID(ctx context.Context, id int64) (*dto.PostResponse, int, error)
	GetPosts(ctx context.Context, page, pageSize int) (*dto.PostsResponse, int, error)
	GetPostsByUserID(ctx context.Context, userID int64, page, pageSize int) (*dto.PostsResponse, int, error)
	UpdatePost(ctx context.Context, userID, postID int64, req dto.UpdatePostRequest) (int, error)
	DeletePost(ctx context.Context, userID, postID int64) (int, error)
}

type postService struct {
	cfg      *config.Config
	postRepo post.PostRepository
	db       *sql.DB
}

func NewService(cfg *config.Config, postRepo post.PostRepository, db *sql.DB) PostService {
	return &postService{
		cfg:      cfg,
		postRepo: postRepo,
		db:       db,
	}
}
