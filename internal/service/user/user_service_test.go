package user

import (
	"context"
	"errors"
	"go-twitter/internal/config"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Mock UserRepository for testing
type mockUserRepository struct {
	getUserByEmailOrUsernameFunc func(ctx context.Context, email, username string) (*model.UserModel, error)
	createUserFunc               func(ctx context.Context, user *model.UserModel) (int64, error)
	getRefreshTokenFunc          func(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error)
	storeRefreshTokenFunc        func(ctx context.Context, model *model.RefreshTokenModel) error
	getUserByIDFunc              func(ctx context.Context, id int64) (*model.UserModel, error)
	getRefreshTokenByTokenFunc   func(ctx context.Context, token string) (*model.RefreshTokenModel, error)
	deleteRefreshTokenFunc       func(ctx context.Context, token string) error
	updateUserFunc               func(ctx context.Context, user *model.UserModel) error
}

func (m *mockUserRepository) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*model.UserModel, error) {
	if m.getUserByEmailOrUsernameFunc != nil {
		return m.getUserByEmailOrUsernameFunc(ctx, email, username)
	}
	return nil, nil
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user *model.UserModel) (int64, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	return 0, nil
}

func (m *mockUserRepository) GetRefreshToken(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error) {
	if m.getRefreshTokenFunc != nil {
		return m.getRefreshTokenFunc(ctx, id, now)
	}
	return nil, nil
}

func (m *mockUserRepository) StoreRefreshToken(ctx context.Context, model *model.RefreshTokenModel) error {
	if m.storeRefreshTokenFunc != nil {
		return m.storeRefreshTokenFunc(ctx, model)
	}
	return nil
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id int64) (*model.UserModel, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepository) GetRefreshTokenByToken(ctx context.Context, token string) (*model.RefreshTokenModel, error) {
	if m.getRefreshTokenByTokenFunc != nil {
		return m.getRefreshTokenByTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *mockUserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if m.deleteRefreshTokenFunc != nil {
		return m.deleteRefreshTokenFunc(ctx, token)
	}
	return nil
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, user *model.UserModel) error {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, user)
	}
	return nil
}

// Test Register
func TestRegister_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, nil // User doesn't exist
		},
		createUserFunc: func(ctx context.Context, user *model.UserModel) (int64, error) {
			return 123, nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	id, status, err := service.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	if id != 123 {
		t.Errorf("Expected user ID 123, got %d", id)
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return &model.UserModel{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			}, nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	id, status, err := service.Register(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists' error, got: %v", err)
	}

	if status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	if id != 0 {
		t.Errorf("Expected user ID 0, got %d", id)
	}
}

func TestRegister_DatabaseErrorOnCheck(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, errors.New("database connection error")
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	id, status, err := service.Register(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if id != 0 {
		t.Errorf("Expected user ID 0, got %d", id)
	}
}

func TestRegister_DatabaseErrorOnCreate(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, nil
		},
		createUserFunc: func(ctx context.Context, user *model.UserModel) (int64, error) {
			return 0, errors.New("failed to create user")
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	id, status, err := service.Register(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if id != 0 {
		t.Errorf("Expected user ID 0, got %d", id)
	}
}

func TestRegister_PasswordIsHashed(t *testing.T) {
	var capturedUser *model.UserModel

	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, nil
		},
		createUserFunc: func(ctx context.Context, user *model.UserModel) (int64, error) {
			capturedUser = user
			return 123, nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	plainPassword := "mySecurePassword123"
	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: plainPassword,
	}

	_, _, err := service.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify password was hashed
	if capturedUser.Password == plainPassword {
		t.Error("Password was not hashed")
	}

	// Verify hash is valid
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.Password), []byte(plainPassword))
	if err != nil {
		t.Errorf("Hashed password doesn't match: %v", err)
	}
}

// Test Login
func TestLogin_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return &model.UserModel{
				ID:       123,
				Username: "testuser",
				Email:    "test@example.com",
				Password: string(hashedPassword),
			}, nil
		},
		getRefreshTokenFunc: func(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error) {
			return nil, nil // No existing refresh token
		},
		storeRefreshTokenFunc: func(ctx context.Context, model *model.RefreshTokenModel) error {
			return nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	accessToken, refreshToken, status, err := service.Login(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	if accessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if refreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, nil // User doesn't exist
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	accessToken, refreshToken, status, err := service.Login(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	if accessToken != "" || refreshToken != "" {
		t.Error("Expected empty tokens")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctPassword"), bcrypt.DefaultCost)

	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return &model.UserModel{
				ID:       123,
				Username: "testuser",
				Email:    "test@example.com",
				Password: string(hashedPassword),
			}, nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongPassword",
	}

	accessToken, refreshToken, status, err := service.Login(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	if accessToken != "" || refreshToken != "" {
		t.Error("Expected empty tokens")
	}
}

func TestLogin_ExistingRefreshToken(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingRefreshToken := "existing-refresh-token"

	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return &model.UserModel{
				ID:       123,
				Username: "testuser",
				Email:    "test@example.com",
				Password: string(hashedPassword),
			}, nil
		},
		getRefreshTokenFunc: func(ctx context.Context, id int64, now time.Time) (*model.RefreshTokenModel, error) {
			return &model.RefreshTokenModel{
				ID:           1,
				UserID:       123,
				RefreshToken: existingRefreshToken,
				ExpiresAt:    time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	accessToken, refreshToken, status, err := service.Login(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	if refreshToken != existingRefreshToken {
		t.Errorf("Expected refresh token %s, got %s", existingRefreshToken, refreshToken)
	}

	if accessToken == "" {
		t.Error("Expected access token, got empty string")
	}
}

func TestLogin_DatabaseError(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (*model.UserModel, error) {
			return nil, errors.New("database error")
		},
	}

	cfg := &config.Config{
		SecreetJwt: "test-secret",
	}

	service := NewService(cfg, mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	_, _, status, err := service.Login(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
}
