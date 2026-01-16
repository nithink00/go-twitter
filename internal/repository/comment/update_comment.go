package comment

import (
	"context"
	"go-twitter/internal/model"
)

func (r *commentRepository) UpdateComment(ctx context.Context, comment *model.CommentModel) error {
	query := `UPDATE comments SET content = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, comment.Content, comment.ID)
	return err
}
