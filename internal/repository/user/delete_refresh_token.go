package user

import (
	"context"
)

func (r *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE refresh_token = ?`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
