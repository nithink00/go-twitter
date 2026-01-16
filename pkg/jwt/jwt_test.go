package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateToken_Success(t *testing.T) {
	// Arrange
	id := int64(123)
	username := "testuser"
	secretKey := "test-secret-key"

	// Act
	token, err := CreateToken(id, username, secretKey)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == "" {
		t.Fatal("Expected token string, got empty string")
	}

	// Verify token can be parsed and contains correct claims
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("Expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Expected claims to be MapClaims")
	}

	// Verify ID claim
	if claimID, ok := claims["id"].(float64); !ok || int64(claimID) != id {
		t.Errorf("Expected id claim to be %d, got %v", id, claims["id"])
	}

	// Verify username claim
	if claimUsername, ok := claims["username"].(string); !ok || claimUsername != username {
		t.Errorf("Expected username claim to be %s, got %v", username, claims["username"])
	}

	// Verify expiration is approximately 60 minutes from now
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		expectedExp := time.Now().Add(60 * time.Minute)
		diff := expTime.Sub(expectedExp).Abs()

		// Allow 5 second margin for test execution time
		if diff > 5*time.Second {
			t.Errorf("Expected expiration to be around 60 minutes from now, got %v", expTime)
		}
	} else {
		t.Fatal("Expected exp claim to be present")
	}
}

func TestCreateToken_DifferentUsers(t *testing.T) {
	// Test that different users get different tokens
	secretKey := "test-secret-key"

	token1, err1 := CreateToken(1, "user1", secretKey)
	token2, err2 := CreateToken(2, "user2", secretKey)

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	if token1 == token2 {
		t.Error("Expected different tokens for different users")
	}
}

func TestCreateToken_VerifySigningMethod(t *testing.T) {
	// Ensure the token uses HS256
	id := int64(456)
	username := "testuser"
	secretKey := "test-secret-key"

	token, _ := CreateToken(id, username, secretKey)

	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Expected HMAC signing method, got %v", token.Method)
		}

		if token.Method.Alg() != "HS256" {
			t.Errorf("Expected HS256 algorithm, got %s", token.Method.Alg())
		}

		return []byte(secretKey), nil
	})

	if parsedToken == nil {
		t.Fatal("Failed to parse token")
	}
}

func TestCreateToken_InvalidSecretKeyVerification(t *testing.T) {
	// Create token with one secret, try to verify with another
	id := int64(789)
	username := "testuser"
	secretKey := "correct-secret"
	wrongSecret := "wrong-secret"

	token, err := CreateToken(id, username, secretKey)
	if err != nil {
		t.Fatalf("Expected no error creating token, got: %v", err)
	}

	// Try to parse with wrong secret
	_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrongSecret), nil
	})

	if err == nil {
		t.Error("Expected error when verifying token with wrong secret, got nil")
	}
}

func TestCreateToken_EmptyValues(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		username  string
		secretKey string
		wantError bool
	}{
		{
			name:      "Empty username",
			id:        1,
			username:  "",
			secretKey: "secret",
			wantError: false, // JWT creation should succeed even with empty username
		},
		{
			name:      "Empty secret key",
			id:        1,
			username:  "user",
			secretKey: "",
			wantError: false, // JWT creation succeeds with empty secret
		},
		{
			name:      "Zero ID",
			id:        0,
			username:  "user",
			secretKey: "secret",
			wantError: false, // Zero ID is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(tt.id, tt.username, tt.secretKey)

			if tt.wantError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if !tt.wantError && token == "" {
				t.Error("Expected token string, got empty")
			}
		})
	}
}

func TestCreateToken_LongValues(t *testing.T) {
	// Test with very long username and secret
	id := int64(999)
	longUsername := string(make([]byte, 1000)) // Very long username
	longSecret := string(make([]byte, 1000))   // Very long secret

	token, err := CreateToken(id, longUsername, longSecret)

	if err != nil {
		t.Fatalf("Expected no error with long values, got: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated with long values")
	}
}

func TestCreateToken_NegativeID(t *testing.T) {
	// Test with negative ID
	id := int64(-123)
	username := "testuser"
	secretKey := "secret"

	token, err := CreateToken(id, username, secretKey)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify the negative ID is preserved
	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	if claimID, ok := claims["id"].(float64); !ok || int64(claimID) != id {
		t.Errorf("Expected id claim to be %d, got %v", id, claims["id"])
	}
}

func TestCreateToken_SpecialCharactersInUsername(t *testing.T) {
	// Test with special characters
	tests := []string{
		"user@example.com",
		"user name with spaces",
		"user-name_123",
		"用户", // Unicode characters
		"user\nwith\nnewlines",
		"user'with\"quotes",
	}

	secretKey := "secret"

	for _, username := range tests {
		t.Run(username, func(t *testing.T) {
			token, err := CreateToken(1, username, secretKey)

			if err != nil {
				t.Errorf("Expected no error for username %q, got: %v", username, err)
			}

			// Verify the username is preserved correctly
			parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			claims, _ := parsedToken.Claims.(jwt.MapClaims)
			if claimUsername, ok := claims["username"].(string); !ok || claimUsername != username {
				t.Errorf("Expected username %q, got %q", username, claimUsername)
			}
		})
	}
}
