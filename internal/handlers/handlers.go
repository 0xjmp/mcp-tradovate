// Package handlers provides request handlers for the MCP (Market Connection Protocol) server.
// It implements handlers for various trading operations and market data requests,
// validating input parameters and coordinating with the Tradovate client.
package handlers

import (
	"fmt"
	"time"

	"github.com/0xjmp/mcp-tradovate/internal/client"
	"github.com/0xjmp/mcp-tradovate/internal/models"
)

// Handler represents a request handler with its description and implementation.
type Handler struct {
	Description string                                            // Human-readable description of the handler's purpose
	Handler     func(map[string]interface{}) (interface{}, error) // Function that processes the request
}

// Handlers is a map of handler names to their implementations.
type Handlers map[string]Handler

// NewHandlers creates a new set of handlers using the provided Tradovate client.
// It initializes all available handlers with their descriptions and implementations.
func NewHandlers(client client.TradovateClientInterface) Handlers {
	return map[string]Handler{
		"authenticate": {
			Description: "Authenticate with Tradovate API",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return handleAuthenticate(client)
			},
		},
		"getAccounts": {
			Description: "Get all accounts for the authenticated user",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetAccounts()
			},
		},
		"getPositions": {
			Description: "Get current positions",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetPositions()
			},
		},
		"placeOrder": {
			Description: "Place a new order",
			Handler:     handlePlaceOrder(client).(func(map[string]interface{}) (interface{}, error)),
		},
		"cancelOrder": {
			Description: "Cancel an existing order",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				orderID := int(params["orderId"].(float64))
				if err := client.CancelOrder(orderID); err != nil {
					return nil, err
				}
				return map[string]bool{"success": true}, nil
			},
		},
		"getFills": {
			Description: "Get fills for a specific order",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				orderID := int(params["orderId"].(float64))
				return client.GetFills(orderID)
			},
		},
		"getContracts": {
			Description: "Get available contracts",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				return client.GetContracts()
			},
		},
		"getMarketData": {
			Description: "Get real-time market data for a contract",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				contractID := int(params["contractId"].(float64))
				return client.GetMarketData(contractID)
			},
		},
		"getHistoricalData": {
			Description: "Get historical price data for a contract",
			Handler:     handleGetHistoricalData(client).(func(map[string]interface{}) (interface{}, error)),
		},
		"setRiskLimits": {
			Description: "Set risk limits for an account",
			Handler:     handleSetRiskLimits(client).(func(map[string]interface{}) (interface{}, error)),
		},
		"getRiskLimits": {
			Description: "Get current risk management limits for an account",
			Handler: func(params map[string]interface{}) (interface{}, error) {
				accountID := int(params["accountId"].(float64))
				return client.GetRiskLimits(accountID)
			},
		},
	}
}

// handleAuthenticate processes authentication requests.
// It calls the Tradovate client's Authenticate method and returns the response.
func handleAuthenticate(client client.TradovateClientInterface) (interface{}, error) {
	return client.Authenticate()
}

// handlePlaceOrder processes order placement requests.
// Required parameters:
// - accountId: (float64) The account ID to place the order for
// - contractId: (float64) The contract ID to trade
// - orderType: (string) The type of order (e.g., "Market", "Limit")
// - quantity: (float64) The number of contracts to trade
// - timeInForce: (string) The time in force for the order
// Optional parameters:
// - price: (float64) The limit price (required for limit orders)
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

// handleSetRiskLimits processes risk limit update requests.
// Required parameters:
// - accountId: (float64) The account ID to set limits for
// - dayMaxLoss: (float64) Maximum loss allowed per day
// - maxDrawdown: (float64) Maximum drawdown allowed
// - maxPositionQty: (float64) Maximum position size allowed
// - trailingStop: (float64) Trailing stop percentage
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

// handleGetHistoricalData processes historical market data requests.
// Required parameters:
// - contractId: (float64) The contract ID to get data for
// - startTime: (string) Start time in RFC3339 format
// - endTime: (string) End time in RFC3339 format
// - interval: (string) Time interval for data points
func handleGetHistoricalData(client client.TradovateClientInterface) interface{} {
	return func(params map[string]interface{}) (interface{}, error) {
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
	}
}

// validateRequiredParams checks if all required parameters are present in the request.
// It returns an error if any required parameter is missing.
func validateRequiredParams(params map[string]interface{}, required []string) error {
	for _, field := range required {
		if _, ok := params[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

// assertFloat64 attempts to convert an interface{} to float64.
// It returns an error if the conversion fails.
func assertFloat64(value interface{}, paramName string) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("invalid type assertion for %s", paramName)
	}
}

// assertString attempts to convert an interface{} to string.
// It returns an error if the conversion fails.
func assertString(value interface{}, paramName string) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid type assertion for %s", paramName)
	}
}
