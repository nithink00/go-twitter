package user

import (
	"context"
	"go-twitter/internal/config"
	"go-twitter/internal/dto"
	"go-twitter/internal/repository/user"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (int64, int, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, string, int, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (string, string, int, error)
	Logout(ctx context.Context, req dto.LogoutRequest) (int, error)
	GetUserByID(ctx context.Context, id int64) (*dto.GetUserResponse, int, error)
}

type userService struct {
	cfg *config.Config
	userRepo user.UserRepository
}

func NewService(cfg *config.Config, userRepo user.UserRepository) UserService {
	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}