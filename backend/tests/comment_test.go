package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// Test 10: Internal notes tidak terlihat oleh requester
func TestInternalNotesVisibility(t *testing.T) {
	cleanTables()

	createTestUser("Requester", "req@test.com", "password123", "requester")
	createTestUser("Agent", "agent@test.com", "password123", "agent")

	tokenReq := loginAndGetToken("req@test.com", "password123")
	tokenAgent := loginAndGetToken("agent@test.com", "password123")

	// Buat tiket
	w := performRequest("POST", "/api/tickets", map[string]string{
		"title":       "Tiket internal notes test",
		"description": "Testing internal notes visibility",
	}, tokenReq)

	resp := parseResponse(w)
	data := resp["data"].(map[string]interface{})
	ticketID := fmt.Sprintf("%.0f", data["id"].(float64))

	// Requester tambah komentar publik
	performRequest("POST", "/api/tickets/"+ticketID+"/comments", map[string]interface{}{
		"body":        "Komentar publik dari requester",
		"is_internal": false,
	}, tokenReq)

	// Agent tambah komentar publik
	performRequest("POST", "/api/tickets/"+ticketID+"/comments", map[string]interface{}{
		"body":        "Komentar publik dari agent",
		"is_internal": false,
	}, tokenAgent)

	// Agent tambah internal note
	performRequest("POST", "/api/tickets/"+ticketID+"/comments", map[string]interface{}{
		"body":        "Ini catatan internal agent",
		"is_internal": true,
	}, tokenAgent)

	// Requester lihat komentar - harus 2 (tanpa internal note)
	w = performRequest("GET", "/api/tickets/"+ticketID+"/comments", nil, tokenReq)
	resp = parseResponse(w)

	comments := resp["data"].([]interface{})
	if len(comments) != 2 {
		t.Errorf("Requester should see 2 comments, got %d", len(comments))
	}

	// Agent lihat komentar - harus 3 (termasuk internal note)
	w = performRequest("GET", "/api/tickets/"+ticketID+"/comments", nil, tokenAgent)
	resp = parseResponse(w)

	comments = resp["data"].([]interface{})
	if len(comments) != 3 {
		t.Errorf("Agent should see 3 comments, got %d", len(comments))
	}

	// Requester coba buat internal note - FORBIDDEN
	w = performRequest("POST", "/api/tickets/"+ticketID+"/comments", map[string]interface{}{
		"body":        "Coba internal note dari requester",
		"is_internal": true,
	}, tokenReq)

	if w.Code != http.StatusForbidden {
		t.Errorf("Requester should not create internal notes, got status %d", w.Code)
	}
}