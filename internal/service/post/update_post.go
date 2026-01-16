package post

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
)

func (s *postService) UpdatePost(ctx context.Context, userID, postID int64, req dto.UpdatePostRequest) (int, error) {
	existingPost, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if existingPost == nil {
		return http.StatusNotFound, nil
	}

	if existingPost.UserID != userID {
		return http.StatusForbidden, nil
	}

	post := &model.PostModel{
		ID:      postID,
		Title:   req.Title,
		Content: req.Content,
	}

	err = s.postRepo.UpdatePost(ctx, post)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
