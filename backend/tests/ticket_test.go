package tests

import (
	"fmt"
	"helpdesk-app/backend/models"
	"net/http"
	"testing"
	"time"
)

// Test 6: Requester tidak dapat melihat tiket milik orang lain
func TestRequesterCannotSeeOtherTickets(t *testing.T) {
	cleanTables()

	createTestUser("User A", "usera@test.com", "password123", "requester")
	userB := createTestUser("User B", "userb@test.com", "password123", "requester")

	tokenA := loginAndGetToken("usera@test.com", "password123")
	tokenB := loginAndGetToken("userb@test.com", "password123")

	// User A buat tiket
	ticketBody := map[string]string{
		"title":       "Tiket milik User A",
		"description": "Deskripsi tiket milik User A saja",
	}
	w := performRequest("POST", "/api/tickets", ticketBody, tokenA)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create ticket, status: %d", w.Code)
	}

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := fmt.Sprintf("%.0f", data["id"].(float64))

	// User B coba akses tiket milik User A -> harus ditolak
	w = performRequest("GET", "/api/tickets/"+ticketID, nil, tokenB)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d (requester should not see other's ticket)", w.Code)
	}

	// User B list tiket -> tiket User A tidak boleh muncul
	w = performRequest("GET", "/api/tickets", nil, tokenB)
	resp = parseResponse(w)
	listData := resp["data"].(map[string]interface{})
	total := listData["total"].(float64)

	if total != 0 {
		t.Errorf("Expected 0 tickets for User B, got %.0f", total)
	}

	// Verifikasi: agent bisa lihat tiket User A
	createTestUser("Agent", "agent@test.com", "password123", "agent")
	tokenAgent := loginAndGetToken("agent@test.com", "password123")

	w = performRequest("GET", "/api/tickets/"+ticketID, nil, tokenAgent)
	if w.Code != http.StatusOK {
		t.Errorf("Agent should be able to see any ticket, got status %d", w.Code)
	}

	_ = userB // suppress unused warning
}

// Test 7: Transisi status valid - open -> in_progress -> resolved -> closed
func TestValidStatusTransitions(t *testing.T) {
	cleanTables()

	createTestUser("Requester", "req@test.com", "password123", "requester")
	createTestUser("Agent", "agent@test.com", "password123", "agent")

	tokenReq := loginAndGetToken("req@test.com", "password123")
	tokenAgent := loginAndGetToken("agent@test.com", "password123")

	// Buat tiket (status: open)
	w := performRequest("POST", "/api/tickets", map[string]string{
		"title":       "Tiket workflow test",
		"description": "Testing valid status transitions",
	}, tokenReq)

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := fmt.Sprintf("%.0f", data["id"].(float64))

	// open -> in_progress (agent)
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "in_progress",
	}, tokenAgent)

	if w.Code != http.StatusOK {
		t.Errorf("open → in_progress should succeed, got status %d", w.Code)
	}

	// in_progress -> resolved (agent) 
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "resolved",
	}, tokenAgent)

	if w.Code != http.StatusOK {
		t.Errorf("in_progress → resolved should succeed, got status %d", w.Code)
	}

	// resolved -> closed (requester)
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "closed",
	}, tokenReq)

	if w.Code != http.StatusOK {
		t.Errorf("resolved → closed should succeed, got status %d", w.Code)
	}
}

// Test 8: Transisi status tidak valid - ditolak
func TestInvalidStatusTransitions(t *testing.T) {
	cleanTables()

	createTestUser("Requester", "req@test.com", "password123", "requester")
	createTestUser("Agent", "agent@test.com", "password123", "agent")

	tokenReq := loginAndGetToken("req@test.com", "password123")
	tokenAgent := loginAndGetToken("agent@test.com", "password123")

	// Buat tiket (status: open)
	w := performRequest("POST", "/api/tickets", map[string]string{
		"title":       "Tiket invalid transition test",
		"description": "Testing invalid status transitions",
	}, tokenReq)

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := fmt.Sprintf("%.0f", data["id"].(float64))

	// open -> resolved (skip step) - INVALID
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "resolved",
	}, tokenAgent)

	if w.Code != http.StatusForbidden {
		t.Errorf("open → resolved should be forbidden, got status %d", w.Code)
	}

	// open -> in_progress (requester) - WRONG ROLE 
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "in_progress",
	}, tokenReq)

	if w.Code != http.StatusForbidden {
		t.Errorf("requester should not be able to transition open → in_progress, got status %d", w.Code)
	}

	// open -> closed (agent) - WRONG ROLE
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]string{
		"status": "closed",
	}, tokenAgent)

	if w.Code != http.StatusForbidden {
		t.Errorf("agent should not be able to transition open → closed, got status %d", w.Code)
	}
}

// Test 9: Requester tidak boleh assign tiket
func TestRequesterCannotAssignTicket(t *testing.T) {
	cleanTables()

	createTestUser("Requester", "req@test.com", "password123", "requester")
	agent := createTestUser("Agent", "agent@test.com", "password123", "agent")

	tokenReq := loginAndGetToken("req@test.com", "password123")

	// Buat tiket
	w := performRequest("POST", "/api/tickets", map[string]string{
		"title":       "Tiket assign test",
		"description": "Requester seharusnya tidak bisa assign",
	}, tokenReq)

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := fmt.Sprintf("%.0f", data["id"].(float64))

	// Requester coba assign ke agent - FORBIDDEN
	w = performRequest("PATCH", "/api/tickets/"+ticketID, map[string]interface{}{
		"assignee_id": agent.ID,
	}, tokenReq)

	if w.Code != http.StatusForbidden {
		t.Errorf("Requester should not be able to assign, got status %d", w.Code)
	}
}

// Test 10: SLA - Tiket urgent yang belum ditangani ditandai terlambat (is_overdue)
func TestTicketSLAOverdue(t *testing.T) {
	cleanTables()

	createTestUser("Requester", "req@test.com", "password123", "requester")
	tokenReq := loginAndGetToken("req@test.com", "password123")

	// Buat tiket urgent
	w := performRequest("POST", "/api/tickets", map[string]string{
		"title":       "Server Database Crash",
		"description": "Server database utama mengalami down total",
		"priority":    "urgent",
	}, tokenReq)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create ticket: %d", w.Code)
	}

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := uint(data["id"].(float64))

	// Baru dibuat -> belum overdue
	w = performRequest("GET", fmt.Sprintf("/api/tickets/%d", ticketID), nil, tokenReq)
	resp = parseResponse(w)
	data = resp["data"].(map[string]interface{})
	if data["is_overdue"].(bool) {
		t.Errorf("New ticket should not be overdue yet")
	}

	// Simulasikan tiket dibuat 3 jam yang lalu (melebihi SLA 2 jam)
	threeHoursAgo := time.Now().Add(-3 * time.Hour)
	testDB.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("created_at", threeHoursAgo)
	
	// Fetch kembali -> harus ditandai is_overdue: true
	w = performRequest("GET", fmt.Sprintf("/api/tickets/%d", ticketID), nil, tokenReq)
	resp = parseResponse(w)
	data = resp["data"].(map[string]interface{})
	if !data["is_overdue"].(bool) {
		t.Errorf("Urgent ticket older than 2 hours should be marked as overdue")
	}
}