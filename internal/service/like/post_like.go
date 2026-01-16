package like

import (
	"context"
	"net/http"
)

func (s *likeService) LikePost(ctx context.Context, userID, postID int64) (int, error) {
	isLiked, err := s.likeRepo.IsPostLiked(ctx, postID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isLiked {
		return http.StatusConflict, nil
	}

	err = s.likeRepo.LikePost(ctx, postID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (s *likeService) UnlikePost(ctx context.Context, userID, postID int64) (int, error) {
	isLiked, err := s.likeRepo.IsPostLiked(ctx, postID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !isLiked {
		return http.StatusNotFound, nil
	}

	err = s.likeRepo.UnlikePost(ctx, postID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *likeService) GetPostLikesCount(ctx context.Context, postID int64) (int, int, error) {
	count, err := s.likeRepo.GetPostLikesCount(ctx, postID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return count, http.StatusOK, nil
}
