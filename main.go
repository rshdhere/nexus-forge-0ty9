package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	hub := newHub()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/messages", hub.handleMessages)
	mux.HandleFunc("/api/stream", hub.handleStream)
	mux.HandleFunc("/api/presence", hub.handlePresence)
	mux.HandleFunc("/", handleIndex)

	addr := "0.0.0.0:" + port
	log.Println("listening on :" + port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

type Message struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type Hub struct {
	mu       sync.RWMutex
	messages []Message
	clients  map[chan []byte]struct{}
	online   map[string]time.Time
}

func newHub() *Hub {
	return &Hub{
		messages: make([]Message, 0, 64),
		clients:  make(map[chan []byte]struct{}),
		online:   make(map[string]time.Time),
	}
}

func (h *Hub) handleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.mu.RLock()
		out := append([]Message(nil), h.messages...)
		h.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	case http.MethodPost:
		var body struct {
			User string `json:"user"`
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		user := sanitizeName(body.User)
		text := sanitizeText(body.Text)
		if user == "" || text == "" {
			http.Error(w, `{"error":"user and text required"}`, http.StatusBadRequest)
			return
		}
		msg := Message{
			ID:        time.Now().UTC().Format("20060102150405.000000000"),
			User:      user,
			Text:      text,
			CreatedAt: time.Now().UTC(),
		}
		h.mu.Lock()
		h.messages = append(h.messages, msg)
		if len(h.messages) > 500 {
			h.messages = h.messages[len(h.messages)-500:]
		}
		h.online[user] = time.Now().UTC()
		h.mu.Unlock()
		h.broadcast(eventPayload("message", msg))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(msg)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Hub) handlePresence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	user := sanitizeName(body.User)
	if user == "" {
		http.Error(w, `{"error":"user required"}`, http.StatusBadRequest)
		return
	}
	h.mu.Lock()
	h.online[user] = time.Now().UTC()
	names := h.onlineNamesLocked()
	h.mu.Unlock()
	h.broadcast(eventPayload("presence", names))
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "online": names})
}

func (h *Hub) handleStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan []byte, 16)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	snapshot := append([]Message(nil), h.messages...)
	online := h.onlineNamesLocked()
	h.mu.Unlock()

	if data, err := json.Marshal(eventPayload("snapshot", map[string]any{
		"messages": snapshot,
		"online":   online,
	})); err == nil {
		_, _ = w.Write([]byte("data: "))
		_, _ = w.Write(data)
		_, _ = w.Write([]byte("\n\n"))
		flusher.Flush()
	}

	notify := r.Context().Done()
	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
		close(ch)
	}()

	for {
		select {
		case <-notify:
			return
		case payload, open := <-ch:
			if !open {
				return
			}
			_, _ = w.Write([]byte("data: "))
			_, _ = w.Write(payload)
			_, _ = w.Write([]byte("\n\n"))
			flusher.Flush()
		}
	}
}

func (h *Hub) broadcast(payload map[string]any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- data:
		default:
		}
	}
}

func (h *Hub) onlineNamesLocked() []string {
	cutoff := time.Now().UTC().Add(-45 * time.Second)
	names := make([]string, 0, len(h.online))
	for name, seen := range h.online {
		if seen.Before(cutoff) {
			delete(h.online, name)
			continue
		}
		names = append(names, name)
	}
	return names
}

func eventPayload(kind string, data any) map[string]any {
	return map[string]any{"type": kind, "data": data}
}

func sanitizeName(s string) string {
	s = trimSpace(s)
	if len(s) > 24 {
		s = s[:24]
	}
	return s
}

func sanitizeText(s string) string {
	s = trimSpace(s)
	if len(s) > 1000 {
		s = s[:1000]
	}
	return s
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
