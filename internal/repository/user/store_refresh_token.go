package user

import (
	"context"
	"go-twitter/internal/model"
)

func (r *userRepository) StoreRefreshToken(ctx context.Context, model *model.RefreshTokenModel) error {
	query := `INSERT INTO refresh_tokens (user_id, refresh_token, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.UserID, model.RefreshToken, model.ExpiresAt, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}