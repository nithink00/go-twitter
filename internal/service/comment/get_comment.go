package comment

import (
	"context"
	"go-twitter/internal/dto"
	"net/http"
)

func (s *commentService) GetCommentByID(ctx context.Context, id int64) (*dto.CommentResponse, int, error) {
	comment, err := s.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if comment == nil {
		return nil, http.StatusNotFound, nil
	}

	user, err := s.userRepo.GetUserByID(ctx, comment.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	likesCount, err := s.commentRepo.GetCommentLikesCount(ctx, comment.ID)
	if err != nil {
		likesCount = 0
	}

	response := &dto.CommentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		UserID:     comment.UserID,
		Username:   user.Username,
		Content:    comment.Content,
		LikesCount: likesCount,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, http.StatusOK, nil
}
