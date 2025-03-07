package handlers

import (
	"fmt"
	"time"

	"github.com/0jxmp/tradovate-mcp/internal/client"
	"github.com/0jxmp/tradovate-mcp/internal/models"
)

// Handler represents an MCP tool handler
type Handler struct {
	Description string
	Parameters  map[string]interface{}
	Handler     func(params map[string]interface{}) (interface{}, error)
}

// NewHandlers creates all MCP tool handlers
func NewHandlers(client client.TradovateClient) map[string]Handler {
	return map[string]Handler{
		"authenticate": {
			Description: "Authenticate with Tradovate API",
			Parameters: map[string]interface{}{
				"name": map[string]string{
					"type":        "string",
					"description": "Username",
				},
				"password": map[string]string{
					"type":        "string",
					"description": "Password",
				},
				"appId": map[string]string{
					"type":        "string",
					"description": "Application ID",
				},
				"appVersion": map[string]string{
					"type":        "string",
					"description": "Application Version",
				},
				"cid": map[string]string{
					"type":        "string",
					"description": "Client ID",
				},
				"sec": map[string]string{
					"type":        "string",
					"description": "Client Secret",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				creds := models.Credentials{
					Name:       params["name"].(string),
					Password:   params["password"].(string),
					AppID:      params["appId"].(string),
					AppVersion: params["appVersion"].(string),
					CID:        params["cid"].(string),
					SEC:        params["sec"].(string),
				}
				return client.Authenticate(creds)
			},
		},
		"getAccounts": {
			Description: "Get all accounts for the authenticated user",
			Parameters:  map[string]interface{}{},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetAccounts()
			},
		},
		"getPositions": {
			Description: "Get current positions",
			Parameters:  map[string]interface{}{},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetPositions()
			},
		},
		"placeOrder": {
			Description: "Place a new order",
			Parameters: map[string]interface{}{
				"accountId": map[string]string{
					"type":        "integer",
					"description": "Account ID",
				},
				"contractId": map[string]string{
					"type":        "integer",
					"description": "Contract ID",
				},
				"orderType": map[string]string{
					"type":        "string",
					"description": "Order type (Market, Limit, etc.)",
				},
				"price": map[string]string{
					"type":        "number",
					"description": "Order price (required for Limit orders)",
				},
				"quantity": map[string]string{
					"type":        "integer",
					"description": "Order quantity",
				},
				"timeInForce": map[string]string{
					"type":        "string",
					"description": "Time in force (Day, GTC, IOC, etc.)",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				order := models.Order{
					AccountID:   int(params["accountId"].(float64)),
					ContractID:  int(params["contractId"].(float64)),
					OrderType:   params["orderType"].(string),
					Price:       params["price"].(float64),
					Quantity:    int(params["quantity"].(float64)),
					TimeInForce: params["timeInForce"].(string),
				}
				return client.PlaceOrder(order)
			},
		},
		"cancelOrder": {
			Description: "Cancel an existing order",
			Parameters: map[string]interface{}{
				"orderId": map[string]string{
					"type":        "integer",
					"description": "Order ID to cancel",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				orderID := int(params["orderId"].(float64))
				err := client.CancelOrder(orderID)
				if err != nil {
					return nil, err
				}
				return map[string]bool{"success": true}, nil
			},
		},
		"getFills": {
			Description: "Get fills for a specific order",
			Parameters: map[string]interface{}{
				"orderId": map[string]string{
					"type":        "integer",
					"description": "Order ID to get fills for",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				orderID := int(params["orderId"].(float64))
				return client.GetFills(orderID)
			},
		},
		"getContracts": {
			Description: "Get available contracts",
			Parameters:  map[string]interface{}{},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetContracts()
			},
		},
		"getMarketData": {
			Description: "Get real-time market data for a contract",
			Parameters: map[string]interface{}{
				"contractId": map[string]string{
					"type":        "integer",
					"description": "Contract ID to get market data for",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				contractID := int(params["contractId"].(float64))
				return client.GetMarketData(contractID)
			},
		},
		"getHistoricalData": {
			Description: "Get historical price data for a contract",
			Parameters: map[string]interface{}{
				"contractId": map[string]string{
					"type":        "integer",
					"description": "Contract ID to get historical data for",
				},
				"startTime": map[string]string{
					"type":        "string",
					"description": "Start time in ISO 8601 format",
				},
				"endTime": map[string]string{
					"type":        "string",
					"description": "End time in ISO 8601 format",
				},
				"interval": map[string]string{
					"type":        "string",
					"description": "Time interval (1m, 5m, 15m, 1h, 1d)",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				startTime, err := time.Parse(time.RFC3339, params["startTime"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid start time: %w", err)
				}

				endTime, err := time.Parse(time.RFC3339, params["endTime"].(string))
				if err != nil {
					return nil, fmt.Errorf("invalid end time: %w", err)
				}

				return client.GetHistoricalData(
					int(params["contractId"].(float64)),
					startTime,
					endTime,
					params["interval"].(string),
				)
			},
		},
		"setRiskLimits": {
			Description: "Set risk management limits for an account",
			Parameters: map[string]interface{}{
				"accountId": map[string]string{
					"type":        "integer",
					"description": "Account ID to set limits for",
				},
				"maxPositions": map[string]string{
					"type":        "integer",
					"description": "Maximum number of positions allowed",
				},
				"maxLoss": map[string]string{
					"type":        "number",
					"description": "Maximum total loss allowed",
				},
				"dailyMaxLoss": map[string]string{
					"type":        "number",
					"description": "Maximum daily loss allowed",
				},
				"marginPercent": map[string]string{
					"type":        "number",
					"description": "Required margin percentage",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				limits := models.RiskLimit{
					AccountID:     int(params["accountId"].(float64)),
					MaxPositions:  int(params["maxPositions"].(float64)),
					MaxLoss:       params["maxLoss"].(float64),
					DailyMaxLoss:  params["dailyMaxLoss"].(float64),
					MarginPercent: params["marginPercent"].(float64),
				}

				err := client.SetRiskLimits(limits)
				if err != nil {
					return nil, err
				}
				return map[string]bool{"success": true}, nil
			},
		},
		"getRiskLimits": {
			Description: "Get current risk management limits for an account",
			Parameters: map[string]interface{}{
				"accountId": map[string]string{
					"type":        "integer",
					"description": "Account ID to get limits for",
				},
			},
			Handler: func(params map[string]interface{}) (interface{}, error) {
				accountID := int(params["accountId"].(float64))
				return client.GetRiskLimits(accountID)
			},
		},
	}
}
