package refreshtoken

import (
	"encoding/hex"
	"testing"
)

func TestGenerateRefreshToken_Success(t *testing.T) {
	// Act
	token, err := GenerateRefreshToken()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == "" {
		t.Fatal("Expected token string, got empty string")
	}

	// Verify token is hex encoded (should be 36 characters for 18 bytes)
	expectedLength := 36 // 18 bytes * 2 (hex encoding)
	if len(token) != expectedLength {
		t.Errorf("Expected token length to be %d, got %d", expectedLength, len(token))
	}

	// Verify token is valid hex
	_, err = hex.DecodeString(token)
	if err != nil {
		t.Errorf("Expected token to be valid hex string, got error: %v", err)
	}
}

func TestGenerateRefreshToken_Uniqueness(t *testing.T) {
	// Generate multiple tokens and verify they are unique
	const numTokens = 100
	tokens := make(map[string]bool)

	for i := 0; i < numTokens; i++ {
		token, err := GenerateRefreshToken()
		if err != nil {
			t.Fatalf("Error generating token %d: %v", i, err)
		}

		if tokens[token] {
			t.Errorf("Duplicate token generated: %s", token)
		}

		tokens[token] = true
	}

	if len(tokens) != numTokens {
		t.Errorf("Expected %d unique tokens, got %d", numTokens, len(tokens))
	}
}

func TestGenerateRefreshToken_Format(t *testing.T) {
	// Generate token and verify format
	token, err := GenerateRefreshToken()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should only contain hex characters (0-9, a-f)
	for _, char := range token {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("Token contains invalid character: %c", char)
		}
	}
}

func TestGenerateRefreshToken_NotEmpty(t *testing.T) {
	// Generate multiple tokens to ensure none are empty
	for i := 0; i < 10; i++ {
		token, err := GenerateRefreshToken()

		if err != nil {
			t.Fatalf("Expected no error on iteration %d, got: %v", i, err)
		}

		if len(token) == 0 {
			t.Errorf("Expected non-empty token on iteration %d", i)
		}
	}
}

func TestGenerateRefreshToken_Randomness(t *testing.T) {
	// Generate two tokens and verify they are different
	token1, err1 := GenerateRefreshToken()
	token2, err2 := GenerateRefreshToken()

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors, got: %v, %v", err1, err2)
	}

	if token1 == token2 {
		t.Error("Expected different tokens from consecutive calls")
	}
}

func TestGenerateRefreshToken_ByteLength(t *testing.T) {
	// Verify the decoded token is 18 bytes
	token, err := GenerateRefreshToken()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	decoded, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("Expected valid hex string, got error: %v", err)
	}

	expectedByteLength := 18
	if len(decoded) != expectedByteLength {
		t.Errorf("Expected decoded token to be %d bytes, got %d", expectedByteLength, len(decoded))
	}
}

func TestGenerateRefreshToken_LowercaseHex(t *testing.T) {
	// Verify that hex encoding is lowercase
	token, err := GenerateRefreshToken()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	for _, char := range token {
		if char >= 'A' && char <= 'F' {
			t.Errorf("Expected lowercase hex, found uppercase character: %c", char)
		}
	}
}

func TestGenerateRefreshToken_ConsistentLength(t *testing.T) {
	// Generate multiple tokens and verify they all have the same length
	const numTokens = 50
	expectedLength := 36

	for i := 0; i < numTokens; i++ {
		token, err := GenerateRefreshToken()

		if err != nil {
			t.Fatalf("Error on iteration %d: %v", i, err)
		}

		if len(token) != expectedLength {
			t.Errorf("Iteration %d: expected length %d, got %d", i, expectedLength, len(token))
		}
	}
}

// Benchmark tests
func BenchmarkGenerateRefreshToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GenerateRefreshToken()
		if err != nil {
			b.Fatalf("Error generating token: %v", err)
		}
	}
}

func BenchmarkGenerateRefreshToken_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GenerateRefreshToken()
			if err != nil {
				b.Fatalf("Error generating token: %v", err)
			}
		}
	})
}
