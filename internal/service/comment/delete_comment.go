package comment

import (
	"context"
	"net/http"
)

func (s *commentService) DeleteComment(ctx context.Context, userID, commentID int64) (int, error) {
	existingComment, err := s.commentRepo.GetCommentByID(ctx, commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if existingComment == nil {
		return http.StatusNotFound, nil
	}

	if existingComment.UserID != userID {
		return http.StatusForbidden, nil
	}

	err = s.commentRepo.DeleteComment(ctx, commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
