package post

import (
	"context"
	"net/http"
)

func (s *postService) DeletePost(ctx context.Context, userID, postID int64) (int, error) {
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

	err = s.postRepo.DeletePost(ctx, postID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
