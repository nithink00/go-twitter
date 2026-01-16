package comment

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
)

func (s *commentService) CreateComment(ctx context.Context, userID, postID int64, req dto.CreateCommentRequest) (int64, int, error) {
	comment := &model.CommentModel{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	commentID, err := s.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return commentID, http.StatusCreated, nil
}
