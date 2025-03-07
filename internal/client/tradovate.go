package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/0jxmp/tradovate-mcp/internal/models"
)

const (
	baseURL = "https://live.tradovate.com/v1"
)

// TradovateClient handles API communication with Tradovate
type TradovateClient struct {
	httpClient  *http.Client
	accessToken string
}

// AuthRequest represents the authentication request body
type AuthRequest struct {
	Name         string `json:"name"`
	Password     string `json:"password"`
	AppID        string `json:"appId"`
	AppVersion   string `json:"appVersion"`
	ClientID     string `json:"cid"`
	ClientSecret string `json:"sec"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken    string `json:"accessToken"`
	MdAccessToken  string `json:"mdAccessToken"`
	ExpirationTime string `json:"expirationTime"`
	UserID         int    `json:"userId"`
	Name           string `json:"name"`
	ErrorText      string `json:"errorText,omitempty"`
}

// NewTradovateClient creates a new Tradovate client
func NewTradovateClient() *TradovateClient {
	return &TradovateClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Authenticate performs the authentication with Tradovate
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

	req, err := http.NewRequest("POST", baseURL+"/auth/accessTokenRequest", bytes.NewBuffer(jsonData))
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

// GetAccessToken returns the current access token
func (c *TradovateClient) GetAccessToken() string {
	return c.accessToken
}

// Account Operations
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

// Trading Operations
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

// Market Data Operations
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

func (c *TradovateClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	return c.httpClient.Do(req)
}
