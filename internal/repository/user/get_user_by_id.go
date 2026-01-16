package user

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
)

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*model.UserModel, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var user model.UserModel
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
