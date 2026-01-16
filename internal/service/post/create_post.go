package post

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
	"time"
)

func (s *postService) CreatePost(ctx context.Context, userID int64, req dto.CreatePostRequest) (int64, int, error) {
	post := &model.PostModel{
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.postRepo.CreatePost(ctx, post)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}
