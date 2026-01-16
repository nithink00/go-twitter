package comment

import (
	"context"
	"go-twitter/internal/dto"
	"math"
	"net/http"
)

func (s *commentService) GetCommentsByPostID(ctx context.Context, postID int64, page, pageSize int) (*dto.CommentsResponse, int, error) {
	offset := (page - 1) * pageSize

	comments, totalCount, err := s.commentRepo.GetCommentsByPostID(ctx, postID, offset, pageSize)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var commentResponses []dto.CommentResponse
	for _, comment := range comments {
		user, err := s.userRepo.GetUserByID(ctx, comment.UserID)
		if err != nil {
			continue
		}

		likesCount, err := s.commentRepo.GetCommentLikesCount(ctx, comment.ID)
		if err != nil {
			likesCount = 0
		}

		commentResponses = append(commentResponses, dto.CommentResponse{
			ID:         comment.ID,
			PostID:     comment.PostID,
			UserID:     comment.UserID,
			Username:   user.Username,
			Content:    comment.Content,
			LikesCount: likesCount,
			CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  comment.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	response := &dto.CommentsResponse{
		Comments:   commentResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, http.StatusOK, nil
}
