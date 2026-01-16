package user

import (
	"context"
	"go-twitter/internal/dto"
	"net/http"
)

func (s *userService) GetUserByID(ctx context.Context, id int64) (*dto.GetUserResponse, int, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if user == nil {
		return nil, http.StatusNotFound, nil
	}

	response := &dto.GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, http.StatusOK, nil
}
