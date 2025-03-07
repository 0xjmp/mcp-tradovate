package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewTradovateClient(t *testing.T) {
	client := NewTradovateClient()
	if client == nil {
		t.Error("Expected non-nil client")
	}
	if client.httpClient == nil {
		t.Error("Expected non-nil HTTP client")
	}
	if client.httpClient.Timeout != 10*time.Second {
		t.Error("Expected 10 second timeout")
	}
}

func TestAuthenticate(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/auth/accessTokenRequest" {
			t.Errorf("Expected /auth/accessTokenRequest path, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content-type, got %s", r.Header.Get("Content-Type"))
		}

		// Decode request body
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		// Verify credentials from environment
		if authReq.Name != os.Getenv("TRADOVATE_USERNAME") {
			t.Errorf("Expected username %s, got %s", os.Getenv("TRADOVATE_USERNAME"), authReq.Name)
		}

		// Send response
		resp := AuthResponse{
			AccessToken:    "test-token",
			MdAccessToken:  "test-md-token",
			ExpirationTime: "2024-12-31T23:59:59Z",
			UserID:         12345,
			Name:           "Test User",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	// Set test environment variables
	os.Setenv("TRADOVATE_USERNAME", "testuser")
	os.Setenv("TRADOVATE_PASSWORD", "testpass")
	os.Setenv("TRADOVATE_APP_ID", "testapp")
	os.Setenv("TRADOVATE_APP_VERSION", "1.0")
	os.Setenv("TRADOVATE_CID", "testcid")
	os.Setenv("TRADOVATE_SEC", "testsec")

	// Test authentication
	client := NewTradovateClient()
	authResp, err := client.Authenticate()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify response
	if authResp.AccessToken != "test-token" {
		t.Errorf("Expected test-token, got %s", authResp.AccessToken)
	}
	if authResp.UserID != 12345 {
		t.Errorf("Expected user ID 12345, got %d", authResp.UserID)
	}
	if client.GetAccessToken() != "test-token" {
		t.Errorf("Expected client to store test-token, got %s", client.GetAccessToken())
	}
}

func TestAuthenticateError(t *testing.T) {
	// Setup test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := AuthResponse{
			ErrorText: "Invalid credentials",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()

	// Test authentication with error
	client := NewTradovateClient()
	_, err := client.Authenticate()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "authentication failed: Invalid credentials" {
		t.Errorf("Expected 'authentication failed: Invalid credentials' error, got %v", err)
	}
}

func TestAuthenticateNetworkError(t *testing.T) {
	// Set invalid base URL to simulate network error
	originalBaseURL := baseURL
	baseURL = "http://invalid-url"
	defer func() { baseURL = originalBaseURL }()

	// Test authentication with network error
	client := NewTradovateClient()
	_, err := client.Authenticate()
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}
