package post

import (
	"context"
	"errors"
	"go-twitter/internal/config"
	"go-twitter/internal/dto"
	"go-twitter/internal/model"
	"net/http"
	"testing"
	"time"
)

// Mock PostRepository for testing
type mockPostRepository struct {
	createPostFunc          func(ctx context.Context, post *model.PostModel) (int64, error)
	getPostByIDFunc         func(ctx context.Context, id int64) (*model.PostModel, error)
	getPostsFunc            func(ctx context.Context, limit, offset int) ([]*model.PostModel, error)
	getPostsByUserIDFunc    func(ctx context.Context, userID int64, limit, offset int) ([]*model.PostModel, error)
	getPostsCountFunc       func(ctx context.Context) (int64, error)
	updatePostFunc          func(ctx context.Context, post *model.PostModel) error
	deletePostFunc          func(ctx context.Context, id int64) error
	getPostWithUserInfoFunc func(ctx context.Context, id int64) (*model.PostModel, string, error)
	getPostsWithUserInfoFunc func(ctx context.Context, limit, offset int) ([]*model.PostModel, []string, error)
}

func (m *mockPostRepository) CreatePost(ctx context.Context, post *model.PostModel) (int64, error) {
	if m.createPostFunc != nil {
		return m.createPostFunc(ctx, post)
	}
	return 0, nil
}

func (m *mockPostRepository) GetPostByID(ctx context.Context, id int64) (*model.PostModel, error) {
	if m.getPostByIDFunc != nil {
		return m.getPostByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockPostRepository) GetPosts(ctx context.Context, limit, offset int) ([]*model.PostModel, error) {
	if m.getPostsFunc != nil {
		return m.getPostsFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockPostRepository) GetPostsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.PostModel, error) {
	if m.getPostsByUserIDFunc != nil {
		return m.getPostsByUserIDFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *mockPostRepository) GetPostsCount(ctx context.Context) (int64, error) {
	if m.getPostsCountFunc != nil {
		return m.getPostsCountFunc(ctx)
	}
	return 0, nil
}

func (m *mockPostRepository) UpdatePost(ctx context.Context, post *model.PostModel) error {
	if m.updatePostFunc != nil {
		return m.updatePostFunc(ctx, post)
	}
	return nil
}

func (m *mockPostRepository) DeletePost(ctx context.Context, id int64) error {
	if m.deletePostFunc != nil {
		return m.deletePostFunc(ctx, id)
	}
	return nil
}

func (m *mockPostRepository) GetPostWithUserInfo(ctx context.Context, id int64) (*model.PostModel, string, error) {
	if m.getPostWithUserInfoFunc != nil {
		return m.getPostWithUserInfoFunc(ctx, id)
	}
	return nil, "", nil
}

func (m *mockPostRepository) GetPostsWithUserInfo(ctx context.Context, limit, offset int) ([]*model.PostModel, []string, error) {
	if m.getPostsWithUserInfoFunc != nil {
		return m.getPostsWithUserInfoFunc(ctx, limit, offset)
	}
	return nil, nil, nil
}

// Test CreatePost
func TestCreatePost_Success(t *testing.T) {
	mockRepo := &mockPostRepository{
		createPostFunc: func(ctx context.Context, post *model.PostModel) (int64, error) {
			return 456, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	id, status, err := service.CreatePost(context.Background(), 123, req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	if id != 456 {
		t.Errorf("Expected post ID 456, got %d", id)
	}
}

func TestCreatePost_DatabaseError(t *testing.T) {
	mockRepo := &mockPostRepository{
		createPostFunc: func(ctx context.Context, post *model.PostModel) (int64, error) {
			return 0, errors.New("database error")
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	id, status, err := service.CreatePost(context.Background(), 123, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}

	if id != 0 {
		t.Errorf("Expected post ID 0, got %d", id)
	}
}

func TestCreatePost_ValidatesTitleAndContent(t *testing.T) {
	var capturedPost *model.PostModel

	mockRepo := &mockPostRepository{
		createPostFunc: func(ctx context.Context, post *model.PostModel) (int64, error) {
			capturedPost = post
			return 1, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	expectedTitle := "Test Title"
	expectedContent := "Test Content"
	expectedUserID := int64(999)

	req := dto.CreatePostRequest{
		Title:   expectedTitle,
		Content: expectedContent,
	}

	service.CreatePost(context.Background(), expectedUserID, req)

	if capturedPost.Title != expectedTitle {
		t.Errorf("Expected title %s, got %s", expectedTitle, capturedPost.Title)
	}

	if capturedPost.Content != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, capturedPost.Content)
	}

	if capturedPost.UserID != expectedUserID {
		t.Errorf("Expected user ID %d, got %d", expectedUserID, capturedPost.UserID)
	}
}

// Test UpdatePost
func TestUpdatePost_Success(t *testing.T) {
	userID := int64(123)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:      postID,
				UserID:  userID,
				Title:   "Old Title",
				Content: "Old Content",
			}, nil
		},
		updatePostFunc: func(ctx context.Context, post *model.PostModel) error {
			return nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.UpdatePostRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	status, err := service.UpdatePost(context.Background(), userID, postID, req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestUpdatePost_PostNotFound(t *testing.T) {
	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return nil, nil // Post not found
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.UpdatePostRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	status, err := service.UpdatePost(context.Background(), 123, 456, req)

	if err != nil {
		t.Errorf("Expected no error for not found, got: %v", err)
	}

	if status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}

func TestUpdatePost_UnauthorizedUser(t *testing.T) {
	ownerID := int64(123)
	differentUserID := int64(999)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:      postID,
				UserID:  ownerID,
				Title:   "Old Title",
				Content: "Old Content",
			}, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.UpdatePostRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	status, err := service.UpdatePost(context.Background(), differentUserID, postID, req)

	if err != nil {
		t.Errorf("Expected no error for forbidden, got: %v", err)
	}

	if status != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, status)
	}
}

func TestUpdatePost_DatabaseError(t *testing.T) {
	userID := int64(123)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return nil, errors.New("database connection error")
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.UpdatePostRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	status, err := service.UpdatePost(context.Background(), userID, postID, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestUpdatePost_UpdateError(t *testing.T) {
	userID := int64(123)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:      postID,
				UserID:  userID,
				Title:   "Old Title",
				Content: "Old Content",
			}, nil
		},
		updatePostFunc: func(ctx context.Context, post *model.PostModel) error {
			return errors.New("update failed")
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.UpdatePostRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	status, err := service.UpdatePost(context.Background(), userID, postID, req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

// Test DeletePost
func TestDeletePost_Success(t *testing.T) {
	userID := int64(123)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:     postID,
				UserID: userID,
			}, nil
		},
		deletePostFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	status, err := service.DeletePost(context.Background(), userID, postID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestDeletePost_PostNotFound(t *testing.T) {
	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return nil, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	status, err := service.DeletePost(context.Background(), 123, 456)

	if err != nil {
		t.Errorf("Expected no error for not found, got: %v", err)
	}

	if status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}

func TestDeletePost_UnauthorizedUser(t *testing.T) {
	ownerID := int64(123)
	differentUserID := int64(999)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:     postID,
				UserID: ownerID,
			}, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	status, err := service.DeletePost(context.Background(), differentUserID, postID)

	if err != nil {
		t.Errorf("Expected no error for forbidden, got: %v", err)
	}

	if status != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, status)
	}
}

func TestDeletePost_DatabaseError(t *testing.T) {
	userID := int64(123)
	postID := int64(456)

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:     postID,
				UserID: userID,
			}, nil
		},
		deletePostFunc: func(ctx context.Context, id int64) error {
			return errors.New("delete failed")
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	status, err := service.DeletePost(context.Background(), userID, postID)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestDeletePost_DeleteCalledWithCorrectID(t *testing.T) {
	userID := int64(123)
	postID := int64(789)
	var deletedID int64

	mockRepo := &mockPostRepository{
		getPostByIDFunc: func(ctx context.Context, id int64) (*model.PostModel, error) {
			return &model.PostModel{
				ID:     postID,
				UserID: userID,
			}, nil
		},
		deletePostFunc: func(ctx context.Context, id int64) error {
			deletedID = id
			return nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	service.DeletePost(context.Background(), userID, postID)

	if deletedID != postID {
		t.Errorf("Expected delete to be called with ID %d, got %d", postID, deletedID)
	}
}

// Test edge cases
func TestCreatePost_SetsTimestamps(t *testing.T) {
	var capturedPost *model.PostModel
	before := time.Now()

	mockRepo := &mockPostRepository{
		createPostFunc: func(ctx context.Context, post *model.PostModel) (int64, error) {
			capturedPost = post
			return 1, nil
		},
	}

	cfg := &config.Config{}
	service := NewService(cfg, mockRepo, nil)

	req := dto.CreatePostRequest{
		Title:   "Test",
		Content: "Content",
	}

	service.CreatePost(context.Background(), 1, req)

	after := time.Now()

	if capturedPost.CreatedAt.Before(before) || capturedPost.CreatedAt.After(after) {
		t.Error("CreatedAt timestamp not set correctly")
	}

	if capturedPost.UpdatedAt.Before(before) || capturedPost.UpdatedAt.After(after) {
		t.Error("UpdatedAt timestamp not set correctly")
	}
}
