package like

import (
	"context"
	"database/sql"
)

func (r *likeRepository) LikePost(ctx context.Context, postID, userID int64) error {
	query := `INSERT INTO post_likes (post_id, user_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, postID, userID)
	return err
}

func (r *likeRepository) UnlikePost(ctx context.Context, postID, userID int64) error {
	query := `DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`
	_, err := r.db.ExecContext(ctx, query, postID, userID)
	return err
}

func (r *likeRepository) IsPostLiked(ctx context.Context, postID, userID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, postID, userID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func (r *likeRepository) GetPostLikesCount(ctx context.Context, postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	return count, err
}
