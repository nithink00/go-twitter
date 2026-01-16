package user

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
	"time"
)

func (r *userRepository) GetRefreshToken(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expires_at FROM refresh_tokens WHERE user_id = ? AND expires_at >= ?`
	row := r.db.QueryRowContext(ctx, query, id, now)
	var result model.RefreshTokenModel
	err := row.Scan(&result.ID, &result.UserID, &result.RefreshToken, &result.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
