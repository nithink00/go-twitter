package comment

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

func (r *commentRepository) GetCommentByID(ctx context.Context, id int64) (*model.CommentModel, error) {
	query := `SELECT id, post_id, user_id, content, deleted_at, created_at, updated_at FROM comments WHERE id = ? AND deleted_at IS NULL`
	row := r.db.QueryRowContext(ctx, query, id)

	var comment model.CommentModel
	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.DeletedAt, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) GetCommentLikesCount(ctx context.Context, commentID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(&count)
	return count, err
}
