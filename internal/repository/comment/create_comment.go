package comment

import (
	"context"
	"go-twitter/internal/model"
)

func (r *commentRepository) CreateComment(ctx context.Context, comment *model.CommentModel) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())`
	result, err := r.db.ExecContext(ctx, query, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
