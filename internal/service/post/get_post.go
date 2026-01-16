package post

import (
	"context"
	"database/sql"
	"go-twitter/internal/dto"
	"net/http"
)

func (s *postService) GetPostByID(ctx context.Context, id int64) (*dto.PostResponse, int, error) {
	post, username, err := s.postRepo.GetPostWithUserInfo(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if post == nil {
		return nil, http.StatusNotFound, nil
	}

	likesCount, err := s.getPostLikesCount(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	commentsCount, err := s.getPostCommentsCount(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	response := &dto.PostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		Username:      username,
		Title:         post.Title,
		Content:       post.Content,
		LikesCount:    likesCount,
		CommentsCount: commentsCount,
		CreatedAt:     post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, http.StatusOK, nil
}

func (s *postService) getPostLikesCount(ctx context.Context, postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = ?`
	var count int
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (s *postService) getPostCommentsCount(ctx context.Context, postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ? AND deleted_at IS NULL`
	var count int
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
