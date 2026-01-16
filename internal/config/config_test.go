package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	// Create a temporary .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `PORT=8080
DATABASE_URL=postgres://localhost:5432/testdb
JWT_SECRET=test-secret-key
DB_HOST=localhost
DB_PORT=3306
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config, got nil")
	}

	// Verify all fields are loaded correctly
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.DBUrlMigration != "postgres://localhost:5432/testdb" {
		t.Errorf("Expected DBUrlMigration to be 'postgres://localhost:5432/testdb', got '%s'", cfg.DBUrlMigration)
	}

	if cfg.SecreetJwt != "test-secret-key" {
		t.Errorf("Expected SecreetJwt to be 'test-secret-key', got '%s'", cfg.SecreetJwt)
	}

	if cfg.DBHost != "localhost" {
		t.Errorf("Expected DBHost to be 'localhost', got '%s'", cfg.DBHost)
	}

	if cfg.DBPort != "3306" {
		t.Errorf("Expected DBPort to be '3306', got '%s'", cfg.DBPort)
	}

	if cfg.DBUser != "testuser" {
		t.Errorf("Expected DBUser to be 'testuser', got '%s'", cfg.DBUser)
	}

	if cfg.DBPassword != "testpass" {
		t.Errorf("Expected DBPassword to be 'testpass', got '%s'", cfg.DBPassword)
	}

	if cfg.DBName != "testdb" {
		t.Errorf("Expected DBName to be 'testdb', got '%s'", cfg.DBName)
	}
}

func TestLoadConfig_MissingEnvFile(t *testing.T) {
	// Create a temp directory without .env file
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err == nil {
		t.Error("Expected error when .env file is missing, got nil")
	}

	if cfg != nil {
		t.Error("Expected nil config when error occurs, got non-nil")
	}
}

func TestLoadConfig_EmptyEnvFile(t *testing.T) {
	// Create an empty .env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(envFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Clear environment variables that might be set
	os.Clearenv()

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert - should succeed but all fields will be empty
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config, got nil")
	}

	// All fields should be empty strings
	if cfg.Port != "" {
		t.Errorf("Expected empty Port, got '%s'", cfg.Port)
	}

	if cfg.DBUrlMigration != "" {
		t.Errorf("Expected empty DBUrlMigration, got '%s'", cfg.DBUrlMigration)
	}

	if cfg.SecreetJwt != "" {
		t.Errorf("Expected empty SecreetJwt, got '%s'", cfg.SecreetJwt)
	}
}

func TestLoadConfig_PartialEnvFile(t *testing.T) {
	// Create a .env file with only some variables
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `PORT=8080
JWT_SECRET=test-secret
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	os.Clearenv()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify present values
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.SecreetJwt != "test-secret" {
		t.Errorf("Expected SecreetJwt to be 'test-secret', got '%s'", cfg.SecreetJwt)
	}

	// Missing values should be empty
	if cfg.DBHost != "" {
		t.Errorf("Expected empty DBHost, got '%s'", cfg.DBHost)
	}

	if cfg.DBUrlMigration != "" {
		t.Errorf("Expected empty DBUrlMigration, got '%s'", cfg.DBUrlMigration)
	}
}

func TestLoadConfig_WithComments(t *testing.T) {
	// Create a .env file with comments
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `# Database configuration
DB_HOST=localhost
DB_PORT=3306

# JWT configuration
JWT_SECRET=secret123

# Server configuration
PORT=8080
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	os.Clearenv()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify values are loaded correctly despite comments
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected DBHost to be 'localhost', got '%s'", cfg.DBHost)
	}

	if cfg.DBPort != "3306" {
		t.Errorf("Expected DBPort to be '3306', got '%s'", cfg.DBPort)
	}

	if cfg.SecreetJwt != "secret123" {
		t.Errorf("Expected SecreetJwt to be 'secret123', got '%s'", cfg.SecreetJwt)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
}

func TestLoadConfig_WithQuotedValues(t *testing.T) {
	// Test with quoted values in .env
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `PORT="8080"
DB_PASSWORD="pass with spaces"
JWT_SECRET='single-quoted-secret'
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	os.Clearenv()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// godotenv should handle quoted values
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}

	if cfg.DBPassword != "pass with spaces" {
		t.Errorf("Expected DBPassword to be 'pass with spaces', got '%s'", cfg.DBPassword)
	}
}

func TestLoadConfig_SpecialCharacters(t *testing.T) {
	// Test with special characters in values
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	envContent := `DB_PASSWORD=p@ssw0rd!#$%
JWT_SECRET=secret-with-dashes_and_underscores
DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	os.Clearenv()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Act
	cfg, err := LoadConfig()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.DBPassword != "p@ssw0rd!#$%" {
		t.Errorf("Expected DBPassword with special chars, got '%s'", cfg.DBPassword)
	}

	if cfg.SecreetJwt != "secret-with-dashes_and_underscores" {
		t.Errorf("Expected SecreetJwt with dashes and underscores, got '%s'", cfg.SecreetJwt)
	}

	if cfg.DBUrlMigration != "postgres://user:pass@localhost:5432/db?sslmode=disable" {
		t.Errorf("Expected DATABASE_URL with special chars, got '%s'", cfg.DBUrlMigration)
	}
}

func TestConfig_StructFields(t *testing.T) {
	// Test Config struct has all expected fields
	cfg := &Config{
		Port:           "test-port",
		DBUrlMigration: "test-url",
		SecreetJwt:     "test-secret",
		DBHost:         "test-host",
		DBPort:         "test-db-port",
		DBUser:         "test-user",
		DBPassword:     "test-password",
		DBName:         "test-db-name",
	}

	// Verify all fields can be set and retrieved
	if cfg.Port != "test-port" {
		t.Errorf("Port field not working correctly")
	}
	if cfg.DBUrlMigration != "test-url" {
		t.Errorf("DBUrlMigration field not working correctly")
	}
	if cfg.SecreetJwt != "test-secret" {
		t.Errorf("SecreetJwt field not working correctly")
	}
	if cfg.DBHost != "test-host" {
		t.Errorf("DBHost field not working correctly")
	}
	if cfg.DBPort != "test-db-port" {
		t.Errorf("DBPort field not working correctly")
	}
	if cfg.DBUser != "test-user" {
		t.Errorf("DBUser field not working correctly")
	}
	if cfg.DBPassword != "test-password" {
		t.Errorf("DBPassword field not working correctly")
	}
	if cfg.DBName != "test-db-name" {
		t.Errorf("DBName field not working correctly")
	}
}
