package comment

import (
	"context"
)

func (r *commentRepository) DeleteComment(ctx context.Context, id int64) error {
	query := `UPDATE comments SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
