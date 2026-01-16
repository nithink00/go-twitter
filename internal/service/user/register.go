package user

import (
	"context"
	"errors"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(ctx context.Context, req dto.RegisterRequest) (int64, int, error) {
	// check user already exists
	userExists, err := s.userRepo.GetUserByEmailOrUsername(ctx, req.Email, req.Username)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	if userExists != nil {
		return 0, http.StatusBadRequest, errors.New("user already exists")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	// create user
	now := time.Now()
	user := &model.UserModel{
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save user
	id, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	return id, http.StatusOK, nil
}
