package user

import (
	"context"
	"errors"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"go-twitter/pkg/jwt"
	"go-twitter/pkg/refreshtoken"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (string, string, int, error) {
	// check user exist
	userExist, err := s.userRepo.GetUserByEmailOrUsername(ctx, req.Email, "")
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	if userExist == nil {
		return "", "", http.StatusBadRequest, errors.New("Please check your credentials and try again")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(req.Password)); err != nil {
		return "", "", http.StatusBadRequest, errors.New("Please check your credentials and try again")
	}
	// generate access token

	accessToken, err := jwt.CreateToken(userExist.ID, userExist.Username, s.cfg.SecreetJwt)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	// get refresh token if exist
	now := time.Now()
	refreshToken, err := s.userRepo.GetRefreshToken(ctx, userExist.ID, now)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	if refreshToken != nil {
		return accessToken, refreshToken.RefreshToken, http.StatusOK, nil
	}
	// generate refresh token
	refreshTokenString, err := refreshtoken.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	// save refresh token
	err = s.userRepo.StoreRefreshToken(ctx, &model.RefreshTokenModel{
		UserID: userExist.ID,
		RefreshToken: refreshTokenString,
		ExpiresAt: now.Add(time.Hour * 24 * 7),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}
	// return access token and refresh token
	return accessToken, refreshTokenString, http.StatusOK, nil
}
