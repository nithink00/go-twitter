package user

import (
	"context"
	"go-twitter/internal/model"
)

func (r *userRepository) CreateUser(ctx context.Context, user *model.UserModel) (int64, error) {
	query := `INSERT INTO users (email, username, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}	