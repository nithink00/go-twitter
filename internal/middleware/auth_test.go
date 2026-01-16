package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Helper function to create a valid JWT token for testing
func createTestToken(userID int, secretKey string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(expiry).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func TestNewAuthMiddleware(t *testing.T) {
	secretKey := "test-secret"

	middleware := NewAuthMiddleware(secretKey)

	if middleware == nil {
		t.Fatal("Expected middleware instance, got nil")
	}

	if middleware.jwtSecret != secretKey {
		t.Errorf("Expected jwtSecret to be %s, got %s", secretKey, middleware.jwtSecret)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	secretKey := "test-secret-key"
	userID := 123

	// Create valid token
	token, err := createTestToken(userID, secretKey, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Setup
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	middleware := NewAuthMiddleware(secretKey)

	// Add middleware and handler
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// Act
	router.ServeHTTP(w, c.Request)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify user_id is set in context (check in handler)
	router.GET("/test2", middleware.RequireAuth(), func(c *gin.Context) {
		value, exists := c.Get("user_id")
		if !exists {
			t.Fatal("Expected user_id to be set in context")
		}

		if id, ok := value.(int); !ok || id != userID {
			t.Errorf("Expected user_id to be %d, got %v", userID, value)
		}
		c.Status(http.StatusOK)
	})

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/test2", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)
}

func TestRequireAuth_MissingAuthorizationHeader(t *testing.T) {
	secretKey := "test-secret-key"

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	// No Authorization header set

	middleware := NewAuthMiddleware(secretKey)

	// Act
	handler := middleware.RequireAuth()
	handler(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Expected request to be aborted")
	}
}

func TestRequireAuth_InvalidHeaderFormat(t *testing.T) {
	secretKey := "test-secret-key"

	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "No Bearer prefix",
			header: "some-token",
		},
		{
			name:   "Wrong prefix",
			header: "Basic some-token",
		},
		{
			name:   "Missing token after Bearer",
			header: "Bearer",
		},
		{
			name:   "Empty Bearer",
			header: "Bearer ",
		},
		{
			name:   "Multiple spaces",
			header: "Bearer  token extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Request.Header.Set("Authorization", tt.header)

			middleware := NewAuthMiddleware(secretKey)
			handler := middleware.RequireAuth()
			handler(c)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}

			if !c.IsAborted() {
				t.Error("Expected request to be aborted")
			}
		})
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	secretKey := "test-secret-key"

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Malformed token",
			token: "invalid.token.here",
		},
		{
			name:  "Random string",
			token: "randomstring",
		},
		{
			name:  "Empty token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Request.Header.Set("Authorization", "Bearer "+tt.token)

			middleware := NewAuthMiddleware(secretKey)
			handler := middleware.RequireAuth()
			handler(c)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}

			if !c.IsAborted() {
				t.Error("Expected request to be aborted")
			}
		})
	}
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	secretKey := "test-secret-key"
	userID := 123

	// Create expired token (expired 1 hour ago)
	token, err := createTestToken(userID, secretKey, -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware := NewAuthMiddleware(secretKey)
	handler := middleware.RequireAuth()
	handler(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Expected request to be aborted")
	}
}

func TestRequireAuth_WrongSecretKey(t *testing.T) {
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	userID := 123

	// Create token with correct secret
	token, err := createTestToken(userID, correctSecret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Try to verify with wrong secret
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	middleware := NewAuthMiddleware(wrongSecret)
	handler := middleware.RequireAuth()
	handler(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Expected request to be aborted")
	}
}

func TestRequireAuth_MissingUserIDClaim(t *testing.T) {
	secretKey := "test-secret-key"

	// Create token without user_id claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secretKey))

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := NewAuthMiddleware(secretKey)
	handler := middleware.RequireAuth()
	handler(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if !c.IsAborted() {
		t.Error("Expected request to be aborted")
	}
}

func TestRequireAuth_DifferentUserIDs(t *testing.T) {
	secretKey := "test-secret-key"

	userIDs := []int{1, 100, 999, 12345}

	for _, userID := range userIDs {
		t.Run("UserID_"+string(rune(userID)), func(t *testing.T) {
			token, err := createTestToken(userID, secretKey, time.Hour)
			if err != nil {
				t.Fatalf("Failed to create test token: %v", err)
			}

			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			middleware := NewAuthMiddleware(secretKey)

			router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
				value, exists := c.Get("user_id")
				if !exists {
					t.Fatal("Expected user_id to be set in context")
				}

				if id, ok := value.(int); !ok || id != userID {
					t.Errorf("Expected user_id to be %d, got %v", userID, value)
				}
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
		})
	}
}

func TestGetUserID_Success(t *testing.T) {
	// Setup
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedID := 123
	c.Set("user_id", expectedID)

	// Act
	id, exists := GetUserID(c)

	// Assert
	if !exists {
		t.Fatal("Expected user_id to exist")
	}

	if id != expectedID {
		t.Errorf("Expected user_id to be %d, got %d", expectedID, id)
	}
}

func TestGetUserID_NotSet(t *testing.T) {
	// Setup
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	// Don't set user_id

	// Act
	id, exists := GetUserID(c)

	// Assert
	if exists {
		t.Error("Expected user_id to not exist")
	}

	if id != 0 {
		t.Errorf("Expected id to be 0 when not found, got %d", id)
	}
}

func TestGetUserID_WrongType(t *testing.T) {
	// Setup
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "not-an-int") // Wrong type

	// Act
	id, exists := GetUserID(c)

	// Assert
	if exists {
		t.Error("Expected exists to be false for wrong type")
	}

	if id != 0 {
		t.Errorf("Expected id to be 0 for wrong type, got %d", id)
	}
}

func TestGetUserID_DifferentValues(t *testing.T) {
	tests := []int{0, 1, -1, 100, 999999}

	for _, expectedID := range tests {
		t.Run("ID_"+string(rune(expectedID)), func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set("user_id", expectedID)

			id, exists := GetUserID(c)

			if !exists {
				t.Fatal("Expected user_id to exist")
			}

			if id != expectedID {
				t.Errorf("Expected user_id to be %d, got %d", expectedID, id)
			}
		})
	}
}

func TestRequireAuth_WrongSigningMethod(t *testing.T) {
	secretKey := "test-secret-key"

	// Create token with RS256 instead of HS256 (requires different key type)
	// For this test, we'll create a token that claims to use a different method
	token := jwt.New(jwt.SigningMethodNone)
	token.Claims = jwt.MapClaims{
		"user_id": float64(123),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := NewAuthMiddleware(secretKey)
	handler := middleware.RequireAuth()
	handler(c)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
