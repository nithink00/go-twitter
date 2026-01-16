package user

import (
	"context"
	"go-twitter/internal/dto"
	"net/http"
)

func (s *userService) Logout(ctx context.Context, req dto.LogoutRequest) (int, error) {
	err := s.userRepo.DeleteRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
