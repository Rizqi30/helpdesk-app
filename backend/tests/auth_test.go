package tests

import (
	"net/http"
	"testing"

	"helpdesk-app/backend/utils"
)

// Test 1: Register berhasil - user baru terdaftar sebagai requester
func TestRegisterSuccess(t *testing.T) {
	cleanTables()

	body := map[string]string{
		"name":     "New User",
		"email":    "newuser@test.com",
		"password": "password123",
	}

	w := performRequest("POST", "/api/auth/register", body, "")
	resp := parseResponse(w)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response should contain 'data' field")
	}

	if data["token"] == nil || data["token"] == "" {
		t.Error("Response should contain a token")
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Response should contain 'user' field")
	}

	if user["role"] != "requester" {
		t.Errorf("Expected role 'requester', got '%s'", user["role"])
	}

	if _, exists := user["password_hash"]; exists {
		t.Error("password_hash should not be in response")
	}
}

// Test 2: Register dengan email duplikat - ditolak
func TestRegisterDuplicateEmail(t *testing.T) {
	cleanTables()
	createTestUser("Existing", "existing@test.com", "password123", "requester")

	body := map[string]string{
		"name":     "Another User",
		"email":    "existing@test.com",
		"password": "password123",
	}

	w := performRequest("POST", "/api/auth/register", body, "")

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}
}

// Test 3: Login berhasil dan JWT valid
func TestLoginAndJWTValidation(t *testing.T) {
	cleanTables()
	createTestUser("Test User", "test@test.com", "password123", "requester")

	// Login
	body := map[string]string{
		"email":    "test@test.com",
		"password": "password123",
	}

	w := performRequest("POST", "/api/auth/login", body, "")
	resp := parseResponse(w)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	data := resp["data"].(map[string]interface{})
	token := data["token"].(string)

	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// Validasi JWT token yang didapat
	jwtSecret := getTestEnv("JWT_SECRET", "test-secret-key")
	claims, err := utils.ValidateToken(token, jwtSecret)
	if err != nil {
		t.Fatalf("Token should be valid, got error: %v", err)
	}

	if claims.Role != "requester" {
		t.Errorf("Expected role 'requester' in token, got '%s'", claims.Role)
	}

	if claims.UserID == 0 {
		t.Error("UserID in token should not be 0")
	}
}

// Test 4: Login dengan password salah - ditolak
func TestLoginWrongPassword(t *testing.T) {
	cleanTables()
	createTestUser("Test User", "test@test.com", "password123", "requester")

	body := map[string]string{
		"email":    "test@test.com",
		"password": "wrongpassword",
	}

	w := performRequest("POST", "/api/auth/login", body, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// Test 5: Akses /api/me tanpa token - ditolak
func TestGetMeWithoutToken(t *testing.T) {
	w := performRequest("GET", "/api/me", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}