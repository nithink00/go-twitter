package like

import (
	"context"
	"database/sql"
)

func (r *likeRepository) LikeComment(ctx context.Context, commentID, userID int64) error {
	query := `INSERT INTO comment_likes (comment_id, user_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, commentID, userID)
	return err
}

func (r *likeRepository) UnlikeComment(ctx context.Context, commentID, userID int64) error {
	query := `DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, commentID, userID)
	return err
}

func (r *likeRepository) IsCommentLiked(ctx context.Context, commentID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, commentID, userID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (r *likeRepository) GetCommentLikesCount(ctx context.Context, commentID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(&count)
	return count, err
}
