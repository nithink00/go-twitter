package post

import (
	"context"
)

func (r *postRepository) DeletePost(ctx context.Context, id int64) error {
	query := `UPDATE posts SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
