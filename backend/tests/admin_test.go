package tests

import (
	"fmt"
	"net/http"
	"testing"

	"helpdesk-app/backend/models"
)

// Test 11: RBAC - hanya admin yang bisa akses /api/users
func TestAdminRBAC(t *testing.T) {
	cleanTables()

	createTestUser("Admin", "admin@test.com", "password123", models.RoleAdmin)
	createTestUser("Agent", "agent@test.com", "password123", models.RoleAgent)
	createTestUser("Requester", "req@test.com", "password123", models.RoleRequester)

	tokenAdmin := loginAndGetToken("admin@test.com", "password123")
	tokenAgent := loginAndGetToken("agent@test.com", "password123")
	tokenReq := loginAndGetToken("req@test.com", "password123")

	// Admin akses /api/users - HARUS BERHASIL (200)
	w := performRequest("GET", "/api/users", nil, tokenAdmin)
	if w.Code != http.StatusOK {
		t.Errorf("Admin should access /api/users, got status %d", w.Code)
	}

	// Agent akses /api/users - HARUS DITOLAK (403)
	w = performRequest("GET", "/api/users", nil, tokenAgent)
	if w.Code != http.StatusForbidden {
		t.Errorf("Agent should be forbidden from /api/users, got status %d", w.Code)
	}

	// Requester akses /api/users - HARUS DITOLAK (403)
	w = performRequest("GET", "/api/users", nil, tokenReq)
	if w.Code != http.StatusForbidden {
		t.Errorf("Requester should be forbidden from /api/users, got status %d", w.Code)
	}

	// Tanpa token - HARUS UNAUTHORIZED (401)
	w = performRequest("GET", "/api/users", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("No token should return 401, got status %d", w.Code)
	}
}

// Test 12: Admin mengubah role user
func TestAdminUpdateRole(t *testing.T) {
	cleanTables()

	createTestUser("Admin", "admin@test.com", "password123", models.RoleAdmin)
	user := createTestUser("Target User", "target@test.com", "password123", models.RoleRequester)

	tokenAdmin := loginAndGetToken("admin@test.com", "password123")

	// Ubah role requester -> agent via PATCH /api/users/:id/role
	updateURL := fmt.Sprintf("/api/users/%d/role", user.ID)
	w := performRequest("PATCH", updateURL, map[string]string{"role": models.RoleAgent}, tokenAdmin)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	resp := parseResponse(w)
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Response data missing or invalid")
	}

	if data["role"] != models.RoleAgent {
		t.Errorf("Expected role 'agent', got '%s'", data["role"])
	}
}