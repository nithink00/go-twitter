package post

import (
	"context"
	"go-twitter/internal/model"
)

func (r *postRepository) CreatePost(ctx context.Context, post *model.PostModel) (int64, error) {
	query := `INSERT INTO posts (user_id, title, content, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())`
	result, err := r.db.ExecContext(ctx, query, post.UserID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
