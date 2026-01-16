package post

import (
	"context"
	"go-twitter/internal/model"
)

func (r *postRepository) UpdatePost(ctx context.Context, post *model.PostModel) error {
	query := `UPDATE posts SET title = ?, content = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	return err
}
