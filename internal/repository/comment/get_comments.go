package comment

import (
	"context"
	"go-twitter/internal/model"
)

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID int64, offset, limit int) ([]*model.CommentModel, int64, error) {
	query := `SELECT id, post_id, user_id, content, deleted_at, created_at, updated_at
	          FROM comments
	          WHERE post_id = ? AND deleted_at IS NULL
	          ORDER BY created_at DESC
	          LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*model.CommentModel
	for rows.Next() {
		var comment model.CommentModel
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.DeletedAt, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, &comment)
	}

	countQuery := `SELECT COUNT(*) FROM comments WHERE post_id = ? AND deleted_at IS NULL`
	var totalCount int64
	err = r.db.QueryRowContext(ctx, countQuery, postID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return comments, totalCount, nil
}
