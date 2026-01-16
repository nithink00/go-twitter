package user

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

func (r *userRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*model.RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expires_at, created_at, updated_at FROM refresh_tokens WHERE refresh_token = ?`
	row := r.db.QueryRowContext(ctx, query, token)

	var result model.RefreshTokenModel
	err := row.Scan(&result.ID, &result.UserID, &result.RefreshToken, &result.ExpiresAt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
