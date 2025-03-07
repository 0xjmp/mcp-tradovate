package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/0xjmp/mcp-tradovate/internal/client"
)

// Request represents an incoming MCP request
type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// Response represents an MCP response
type Response struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var tradovateClient *client.TradovateClient

func init() {
	tradovateClient = client.NewTradovateClient()
}

func main() {
	// Initialize scanner for STDIN
	scanner := bufio.NewScanner(os.Stdin)

	// Process incoming requests
	for scanner.Scan() {
		line := scanner.Text()

		// Parse request
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(req.ID, 400, fmt.Sprintf("Invalid request: %v", err))
			continue
		}

		// Handle request
		switch req.Method {
		case "ping":
			sendResponse(req.ID, "pong")
		case "authenticate":
			handleAuthenticate(req.ID)
		default:
			sendError(req.ID, 404, fmt.Sprintf("Unknown method: %s", req.Method))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading standard input: %v", err)
	}
}

func handleAuthenticate(reqID string) {
	authResp, err := tradovateClient.Authenticate()
	if err != nil {
		sendError(reqID, 401, fmt.Sprintf("Authentication failed: %v", err))
		return
	}

	sendResponse(reqID, map[string]interface{}{
		"status":         "authenticated",
		"token":          authResp.AccessToken,
		"mdToken":        authResp.MdAccessToken,
		"userId":         authResp.UserID,
		"name":           authResp.Name,
		"expirationTime": authResp.ExpirationTime,
	})
}

func sendResponse(id string, result interface{}) {
	resp := Response{
		ID:     id,
		Result: result,
	}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func sendError(id string, code int, message string) {
	resp := Response{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
