package user

import (
	"context"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"go-twitter/pkg/jwt"
	"go-twitter/pkg/refreshtoken"
	"net/http"
	"time"
)

func (s *userService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (string, string, int, error) {
	refreshTokenModel, err := s.userRepo.GetRefreshTokenByToken(ctx, req.RefreshToken)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if refreshTokenModel == nil {
		return "", "", http.StatusUnauthorized, nil
	}

	if refreshTokenModel.ExpiresAt.Before(time.Now()) {
		return "", "", http.StatusUnauthorized, nil
	}

	user, err := s.userRepo.GetUserByID(ctx, refreshTokenModel.UserID)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if user == nil {
		return "", "", http.StatusUnauthorized, nil
	}

	token, err := jwt.CreateToken(user.ID, user.Username, s.cfg.SecreetJwt)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	newRefreshToken, err := refreshtoken.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	err = s.userRepo.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshTokenModelNew := &model.RefreshTokenModel{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    refreshTokenExpiry,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.StoreRefreshToken(ctx, refreshTokenModelNew)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	return token, newRefreshToken, http.StatusOK, nil
}
