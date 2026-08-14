package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handleHealth(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var body map[string]bool
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if !body["ok"] {
		t.Fatalf("expected ok true, got %#v", body)
	}
}

func TestPostMessage(t *testing.T) {
	hub := newHub()
	body := strings.NewReader(`{"user":"Ada","text":"hello nexus"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/messages", body)
	rr := httptest.NewRecorder()
	hub.handleMessages(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	var msg Message
	if err := json.NewDecoder(rr.Body).Decode(&msg); err != nil {
		t.Fatal(err)
	}
	if msg.User != "Ada" || msg.Text != "hello nexus" || msg.ID == "" {
		t.Fatalf("unexpected message %#v", msg)
	}
}

func TestSanitize(t *testing.T) {
	if got := sanitizeName("  longname-that-exceeds-limit-xxx  "); len(got) > 24 {
		t.Fatalf("name too long: %q", got)
	}
	if got := sanitizeText("  hi  "); got != "hi" {
		t.Fatalf("text trim failed: %q", got)
	}
}
