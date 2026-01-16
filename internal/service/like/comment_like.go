package like

import (
	"context"
	"net/http"
)

func (s *likeService) LikeComment(ctx context.Context, userID, commentID int64) (int, error) {
	isLiked, err := s.likeRepo.IsCommentLiked(ctx, commentID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isLiked {
		return http.StatusConflict, nil
	}

	err = s.likeRepo.LikeComment(ctx, commentID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (s *likeService) UnlikeComment(ctx context.Context, userID, commentID int64) (int, error) {
	isLiked, err := s.likeRepo.IsCommentLiked(ctx, commentID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !isLiked {
		return http.StatusNotFound, nil
	}

	err = s.likeRepo.UnlikeComment(ctx, commentID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *likeService) GetCommentLikesCount(ctx context.Context, commentID int64) (int, int, error) {
	count, err := s.likeRepo.GetCommentLikesCount(ctx, commentID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return count, http.StatusOK, nil
}
