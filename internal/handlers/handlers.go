package handlers

import (
	"fmt"
	"time"

	"github.com/0xjmp/mcp-tradovate/internal/client"
	"github.com/0xjmp/mcp-tradovate/internal/models"
)

// Handler represents an MCP tool handler
type Handler struct {
	Description string
	Parameters  interface{}
	Handler     interface{}
}

// NewHandlers creates a new set of MCP tool handlers
func NewHandlers(client client.TradovateClientInterface) map[string]Handler {
	return map[string]Handler{
		"authenticate": {
			Description: "Authenticate with Tradovate API",
			Parameters:  nil,
			Handler:     handleAuthenticate(client),
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
			Handler: handlePlaceOrder(client),
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
			Description: "Set risk limits for an account",
			Parameters: map[string]interface{}{
				"accountId": map[string]string{
					"type":        "integer",
					"description": "Account ID",
				},
				"dayMaxLoss": map[string]string{
					"type":        "number",
					"description": "Maximum daily loss limit",
				},
				"maxDrawdown": map[string]string{
					"type":        "number",
					"description": "Maximum drawdown limit",
				},
				"maxPositionQty": map[string]string{
					"type":        "integer",
					"description": "Maximum position quantity",
				},
				"trailingStop": map[string]string{
					"type":        "number",
					"description": "Trailing stop percentage",
				},
			},
			Handler: handleSetRiskLimits(client),
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

func handleAuthenticate(client client.TradovateClientInterface) interface{} {
	return func() (interface{}, error) {
		return client.Authenticate()
	}
}

func handlePlaceOrder(client client.TradovateClientInterface) interface{} {
	return func(params map[string]interface{}) (interface{}, error) {
		// Validate required fields
		requiredFields := []string{"accountId", "contractId", "orderType", "quantity", "timeInForce"}
		for _, field := range requiredFields {
			if _, ok := params[field]; !ok {
				return nil, fmt.Errorf("missing required field: %s", field)
			}
		}

		// Type assertions with validation
		accountID, ok := params["accountId"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type assertion for accountId")
		}

		contractID, ok := params["contractId"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type assertion for contractId")
		}

		orderType, ok := params["orderType"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid type assertion for orderType")
		}

		quantity, ok := params["quantity"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type assertion for quantity")
		}

		timeInForce, ok := params["timeInForce"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid type assertion for timeInForce")
		}

		// Price is optional for market orders
		var price float64
		if orderType == "Limit" {
			priceVal, ok := params["price"].(float64)
			if !ok {
				return nil, fmt.Errorf("price is required for Limit orders")
			}
			price = priceVal
		}

		order := models.Order{
			AccountID:   int(accountID),
			ContractID:  int(contractID),
			OrderType:   orderType,
			Price:       price,
			Quantity:    int(quantity),
			TimeInForce: timeInForce,
		}

		return client.PlaceOrder(order)
	}
}

func handleSetRiskLimits(client client.TradovateClientInterface) interface{} {
	return func(params map[string]interface{}) (interface{}, error) {
		accountID, ok := params["accountId"].(float64)
		if !ok {
			return nil, fmt.Errorf("missing or invalid accountId")
		}

		dayMaxLoss, ok := params["dayMaxLoss"].(float64)
		if !ok || dayMaxLoss < 0 {
			return nil, fmt.Errorf("missing or invalid dayMaxLoss")
		}

		maxDrawdown, ok := params["maxDrawdown"].(float64)
		if !ok || maxDrawdown < 0 {
			return nil, fmt.Errorf("missing or invalid maxDrawdown")
		}

		maxPositionQty, ok := params["maxPositionQty"].(float64)
		if !ok || maxPositionQty < 0 {
			return nil, fmt.Errorf("missing or invalid maxPositionQty")
		}

		trailingStop, ok := params["trailingStop"].(float64)
		if !ok || trailingStop < 0 {
			return nil, fmt.Errorf("missing or invalid trailingStop")
		}

		limits := models.RiskLimit{
			AccountID:      int(accountID),
			DayMaxLoss:     dayMaxLoss,
			MaxDrawdown:    maxDrawdown,
			MaxPositionQty: int(maxPositionQty),
			TrailingStop:   trailingStop,
		}
		return nil, client.SetRiskLimits(limits)
	}
}
