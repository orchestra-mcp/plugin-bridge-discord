package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWorkspaceClient_ChatSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Response: "hello"})
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "test-token")
	got, err := client.Chat("ws-123", "hi")
	if err != nil {
		t.Fatalf("Chat returned unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("Chat response = %q, want %q", got, "hello")
	}
}

func TestWorkspaceClient_ChatError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "test-token")
	_, err := client.Chat("ws-123", "hi")
	if err == nil {
		t.Fatal("Chat should return error for 500 status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status code 500, got: %v", err)
	}
}

func TestWorkspaceClient_ChatAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Error: "not found"})
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "test-token")
	_, err := client.Chat("ws-123", "hi")
	if err == nil {
		t.Fatal("Chat should return error when API returns error field")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should contain %q, got: %v", "not found", err)
	}
}

func TestWorkspaceClient_AuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Response: "ok"})
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "my-secret-token")
	_, err := client.Chat("ws-123", "hi")
	if err != nil {
		t.Fatalf("Chat returned unexpected error: %v", err)
	}
	want := "Bearer my-secret-token"
	if gotAuth != want {
		t.Errorf("Authorization header = %q, want %q", gotAuth, want)
	}
}

func TestWorkspaceClient_NoToken(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Response: "ok"})
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "")
	_, err := client.Chat("ws-123", "hi")
	if err != nil {
		t.Fatalf("Chat returned unexpected error: %v", err)
	}
	if gotAuth != "" {
		t.Errorf("Authorization header should be empty when no token, got %q", gotAuth)
	}
}

func TestWorkspaceClient_URLConstruction(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Response: "ok"})
	}))
	defer srv.Close()

	client := NewWorkspaceClient(srv.URL, "token")
	_, err := client.Chat("workspace-abc", "hi")
	if err != nil {
		t.Fatalf("Chat returned unexpected error: %v", err)
	}
	want := "/api/repos/workspace-abc/chat"
	if gotPath != want {
		t.Errorf("request path = %q, want %q", gotPath, want)
	}
}
