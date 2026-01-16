package comment

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
)

func (s *commentService) UpdateComment(ctx context.Context, userID, commentID int64, req dto.UpdateCommentRequest) (int, error) {
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

	comment := &model.CommentModel{
		ID:      commentID,
		Content: req.Content,
	}

	err = s.commentRepo.UpdateComment(ctx, comment)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
