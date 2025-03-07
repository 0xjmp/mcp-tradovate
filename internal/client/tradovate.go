// Package client provides a Go client for interacting with the Tradovate API.
// It handles authentication, trading operations, and market data retrieval through
// a clean, type-safe interface.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/0xjmp/mcp-tradovate/internal/models"
)

// TradovateClientInterface defines the interface for Tradovate client operations.
// This interface allows for easy mocking in tests and provides a clear contract
// for implementing alternative client implementations.
type TradovateClientInterface interface {
	// Authenticate performs the initial authentication with Tradovate and returns the auth response.
	Authenticate() (*AuthResponse, error)
	// GetAccounts retrieves all accounts associated with the authenticated user.
	GetAccounts() ([]models.Account, error)
	// GetRiskLimits retrieves the risk limits for a specific account.
	GetRiskLimits(accountID int) (*models.RiskLimit, error)
	// SetRiskLimits updates the risk limits for a specific account.
	SetRiskLimits(limits models.RiskLimit) error
	// PlaceOrder submits a new order to Tradovate.
	PlaceOrder(order models.Order) (*models.Order, error)
	// CancelOrder cancels an existing order by its ID.
	CancelOrder(orderID int) error
	// GetFills retrieves all fills for a specific order.
	GetFills(orderID int) ([]models.Fill, error)
	// GetPositions retrieves all current positions for the authenticated user.
	GetPositions() ([]models.Position, error)
	// GetContracts retrieves all available trading contracts.
	GetContracts() ([]models.Contract, error)
	// GetMarketData retrieves current market data for a specific contract.
	GetMarketData(contractID int) (*models.MarketData, error)
	// GetHistoricalData retrieves historical market data for a specific contract.
	GetHistoricalData(contractID int, startTime, endTime time.Time, interval string) ([]models.HistoricalData, error)
}

// TradovateClient handles API communication with Tradovate.
// It implements the TradovateClientInterface and manages the HTTP client,
// authentication state, and base URL configuration.
type TradovateClient struct {
	httpClient  *http.Client
	accessToken string
	baseURL     string
}

// AuthRequest represents the authentication request body sent to Tradovate.
// All fields are required for successful authentication.
type AuthRequest struct {
	Name         string `json:"name"`       // Username for Tradovate account
	Password     string `json:"password"`   // Password for Tradovate account
	AppID        string `json:"appId"`      // Application ID provided by Tradovate
	AppVersion   string `json:"appVersion"` // Application version string
	ClientID     string `json:"cid"`        // OAuth client ID
	ClientSecret string `json:"sec"`        // OAuth client secret
}

// AuthResponse represents the authentication response from Tradovate.
// A successful response includes tokens and user information.
type AuthResponse struct {
	AccessToken    string `json:"accessToken"`         // JWT token for API access
	MdAccessToken  string `json:"mdAccessToken"`       // JWT token for market data access
	ExpirationTime string `json:"expirationTime"`      // Token expiration time in ISO format
	UserID         int    `json:"userId"`              // Unique identifier for the user
	Name           string `json:"name"`                // Username of the authenticated user
	ErrorText      string `json:"errorText,omitempty"` // Error message if authentication fails
}

// NewTradovateClient creates a new Tradovate client with default configuration.
// It sets up an HTTP client with a 10-second timeout and uses the live Tradovate API URL.
func NewTradovateClient() *TradovateClient {
	return &TradovateClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://live.tradovate.com/v1",
	}
}

// SetBaseURL sets the base URL for API requests.
// This is useful for testing or switching between demo and live environments.
func (c *TradovateClient) SetBaseURL(url string) {
	c.baseURL = url
}

// Authenticate performs the authentication with Tradovate using environment variables.
// Required environment variables:
// - TRADOVATE_USERNAME: Tradovate account username
// - TRADOVATE_PASSWORD: Tradovate account password
// - TRADOVATE_APP_ID: Application ID from Tradovate
// - TRADOVATE_APP_VERSION: Application version string
// - TRADOVATE_CID: OAuth client ID
// - TRADOVATE_SEC: OAuth client secret
func (c *TradovateClient) Authenticate() (*AuthResponse, error) {
	authReq := AuthRequest{
		Name:         os.Getenv("TRADOVATE_USERNAME"),
		Password:     os.Getenv("TRADOVATE_PASSWORD"),
		AppID:        os.Getenv("TRADOVATE_APP_ID"),
		AppVersion:   os.Getenv("TRADOVATE_APP_VERSION"),
		ClientID:     os.Getenv("TRADOVATE_CID"),
		ClientSecret: os.Getenv("TRADOVATE_SEC"),
	}

	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/auth/accessTokenRequest", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if authResp.ErrorText != "" {
		return nil, fmt.Errorf("authentication failed: %s", authResp.ErrorText)
	}

	c.accessToken = authResp.AccessToken
	return &authResp, nil
}

// GetAccessToken returns the current access token.
// This token is used for authenticating subsequent API requests.
func (c *TradovateClient) GetAccessToken() string {
	return c.accessToken
}

// GetAccounts retrieves all accounts associated with the authenticated user.
// Returns a slice of Account objects containing account details and balances.
func (c *TradovateClient) GetAccounts() ([]models.Account, error) {
	resp, err := c.doRequest("GET", "/account/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accounts []models.Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("error decoding accounts: %w", err)
	}

	return accounts, nil
}

// GetRiskLimits retrieves the risk limits for a specific account.
// Parameters:
// - accountID: The unique identifier of the account
func (c *TradovateClient) GetRiskLimits(accountID int) (*models.RiskLimit, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/account/riskLimits/%d", accountID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var limits models.RiskLimit
	if err := json.NewDecoder(resp.Body).Decode(&limits); err != nil {
		return nil, fmt.Errorf("error decoding risk limits: %w", err)
	}

	return &limits, nil
}

// SetRiskLimits updates the risk limits for a specific account.
// The limits parameter must include all required risk limit fields.
func (c *TradovateClient) SetRiskLimits(limits models.RiskLimit) error {
	resp, err := c.doRequest("POST", "/account/setRiskLimits", limits)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set risk limits: status %d", resp.StatusCode)
	}

	return nil
}

// PlaceOrder submits a new order to Tradovate.
// The order parameter must include all required order fields such as
// account ID, contract ID, order type, quantity, and time in force.
func (c *TradovateClient) PlaceOrder(order models.Order) (*models.Order, error) {
	resp, err := c.doRequest("POST", "/order/placeOrder", order)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var placedOrder models.Order
	if err := json.NewDecoder(resp.Body).Decode(&placedOrder); err != nil {
		return nil, fmt.Errorf("error decoding order response: %w", err)
	}

	return &placedOrder, nil
}

// CancelOrder cancels an existing order by its ID.
// Returns an error if the order cannot be cancelled or doesn't exist.
func (c *TradovateClient) CancelOrder(orderID int) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/order/cancel/%d", orderID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to cancel order: status %d", resp.StatusCode)
	}

	return nil
}

// GetFills retrieves all fills for a specific order.
// Parameters:
// - orderID: The unique identifier of the order
func (c *TradovateClient) GetFills(orderID int) ([]models.Fill, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/fill/list/%d", orderID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fills []models.Fill
	if err := json.NewDecoder(resp.Body).Decode(&fills); err != nil {
		return nil, fmt.Errorf("error decoding fills: %w", err)
	}

	return fills, nil
}

// GetPositions retrieves all current positions for the authenticated user.
// Returns a slice of Position objects containing position details and P&L information.
func (c *TradovateClient) GetPositions() ([]models.Position, error) {
	resp, err := c.doRequest("GET", "/position/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var positions []models.Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, fmt.Errorf("error decoding positions: %w", err)
	}

	return positions, nil
}

// GetContracts retrieves all available trading contracts.
// Returns a slice of Contract objects containing contract specifications.
func (c *TradovateClient) GetContracts() ([]models.Contract, error) {
	resp, err := c.doRequest("GET", "/contract/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contracts []models.Contract
	if err := json.NewDecoder(resp.Body).Decode(&contracts); err != nil {
		return nil, fmt.Errorf("error decoding contracts: %w", err)
	}

	return contracts, nil
}

// GetMarketData retrieves current market data for a specific contract.
// Parameters:
// - contractID: The unique identifier of the contract
func (c *TradovateClient) GetMarketData(contractID int) (*models.MarketData, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/md/getQuote/%d", contractID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var marketData models.MarketData
	if err := json.NewDecoder(resp.Body).Decode(&marketData); err != nil {
		return nil, fmt.Errorf("error decoding market data: %w", err)
	}

	return &marketData, nil
}

// GetHistoricalData retrieves historical market data for a specific contract.
// Parameters:
// - contractID: The unique identifier of the contract
// - startTime: The start time for historical data
// - endTime: The end time for historical data
// - interval: The time interval for data points (e.g., "1m", "5m", "1h")
func (c *TradovateClient) GetHistoricalData(contractID int, startTime, endTime time.Time, interval string) ([]models.HistoricalData, error) {
	params := map[string]interface{}{
		"contractId": contractID,
		"startTime":  startTime.Unix(),
		"endTime":    endTime.Unix(),
		"interval":   interval,
	}

	resp, err := c.doRequest("GET", "/md/historical", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []models.HistoricalData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding historical data: %w", err)
	}

	return data, nil
}

// doRequest performs an HTTP request to the Tradovate API.
// It handles request creation, authentication, and error responses.
// Parameters:
// - method: HTTP method (GET, POST, etc.)
// - endpoint: API endpoint path
// - body: Optional request body for POST/PUT requests
func (c *TradovateClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			ErrorText string `json:"errorText"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("status %d", resp.StatusCode)
		}
		resp.Body.Close()
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, errResp.ErrorText)
	}

	return resp, nil
}
