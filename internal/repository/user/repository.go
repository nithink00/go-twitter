package user

import (
	"context"
	"database/sql"
	"go-twitter/internal/model"
	"time"
)


type UserRepository interface {
	GetUserByEmailOrUsername(ctx context.Context, email, username string) (*model.UserModel, error)
	CreateUser(ctx context.Context, user *model.UserModel) (int64, error)
	GetRefreshToken(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error)
	StoreRefreshToken(ctx context.Context, model *model.RefreshTokenModel) error
	GetUserByID(ctx context.Context, id int64) (*model.UserModel, error)
	GetRefreshTokenByToken(ctx context.Context, token string) (*model.RefreshTokenModel, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	UpdateUser(ctx context.Context, user *model.UserModel) error
}

type userRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}