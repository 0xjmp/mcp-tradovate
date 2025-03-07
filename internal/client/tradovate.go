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

const baseURL = "https://live.tradovate.com/v1"

// tradovateClient implements the TradovateClient interface
type tradovateClient struct {
	accessToken string
	client      *http.Client
}

// NewTradovateClient creates a new instance of TradovateClient
func NewTradovateClient() TradovateClient {
	return &tradovateClient{
		client: &http.Client{},
	}
}

// NewTradovateClientWithEnvCredentials creates a new client and automatically authenticates using environment variables
func NewTradovateClientWithEnvCredentials() (TradovateClient, error) {
	client := &tradovateClient{
		client: &http.Client{},
	}

	// Get credentials from environment variables
	creds := models.Credentials{
		Name:       os.Getenv("TRADOVATE_USERNAME"),
		Password:   os.Getenv("TRADOVATE_PASSWORD"),
		AppID:      os.Getenv("TRADOVATE_APP_ID"),
		AppVersion: os.Getenv("TRADOVATE_APP_VERSION"),
		CID:        os.Getenv("TRADOVATE_CID"),
		SEC:        os.Getenv("TRADOVATE_SEC"),
	}

	// Validate required credentials
	if creds.Name == "" || creds.Password == "" || creds.AppID == "" ||
		creds.AppVersion == "" || creds.CID == "" || creds.SEC == "" {
		return nil, fmt.Errorf("missing required environment variables for authentication")
	}

	// Authenticate
	_, err := client.Authenticate(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with environment credentials: %w", err)
	}

	return client, nil
}

func (c *tradovateClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
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

	return c.client.Do(req)
}

// Authentication
func (c *tradovateClient) Authenticate(creds models.Credentials) (*models.AccessToken, error) {
	resp, err := c.doRequest("POST", "/auth/accessTokenRequest", creds)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	var token models.AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	c.accessToken = token.AccessToken
	return &token, nil
}

// Account Operations
func (c *tradovateClient) GetAccounts() ([]models.Account, error) {
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

func (c *tradovateClient) GetRiskLimits(accountID int) (*models.RiskLimit, error) {
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

func (c *tradovateClient) SetRiskLimits(limits models.RiskLimit) error {
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
func (c *tradovateClient) PlaceOrder(order models.Order) (*models.Order, error) {
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

func (c *tradovateClient) CancelOrder(orderID int) error {
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

func (c *tradovateClient) GetFills(orderID int) ([]models.Fill, error) {
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

func (c *tradovateClient) GetPositions() ([]models.Position, error) {
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
func (c *tradovateClient) GetContracts() ([]models.Contract, error) {
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

func (c *tradovateClient) GetMarketData(contractID int) (*models.MarketData, error) {
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

func (c *tradovateClient) GetHistoricalData(contractID int, startTime, endTime time.Time, interval string) ([]models.HistoricalData, error) {
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
