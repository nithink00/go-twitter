package user

import (
	"context"
	"go-twitter/internal/model"
)

func (r *userRepository) UpdateUser(ctx context.Context, user *model.UserModel) error {
	query := `UPDATE users SET username = ?, email = ?, password = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	return err
}
