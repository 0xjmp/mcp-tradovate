package handlers

import (
	"errors"
	"testing"
	"time"

	"github.com/0xjmp/mcp-tradovate/internal/client"
	"github.com/0xjmp/mcp-tradovate/internal/models"
	"github.com/stretchr/testify/assert"
)

// MockTradovateClient is a mock implementation for testing
type MockTradovateClient struct {
	setRiskLimitsFunc     func(models.RiskLimit) error
	authenticateFunc      func() (*client.AuthResponse, error)
	getAccountsFunc       func() ([]models.Account, error)
	placeOrderFunc        func(models.Order) (*models.Order, error)
	cancelOrderFunc       func(int) error
	getFillsFunc          func(int) ([]models.Fill, error)
	getPositionsFunc      func() ([]models.Position, error)
	getContractsFunc      func() ([]models.Contract, error)
	getMarketDataFunc     func(int) (*models.MarketData, error)
	getRiskLimitsFunc     func(int) (*models.RiskLimit, error)
	getHistoricalDataFunc func(int, time.Time, time.Time, string) ([]models.HistoricalData, error)
}

func (m *MockTradovateClient) SetRiskLimits(limits models.RiskLimit) error {
	if m.setRiskLimitsFunc != nil {
		return m.setRiskLimitsFunc(limits)
	}
	return nil
}

func (m *MockTradovateClient) Authenticate() (*client.AuthResponse, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc()
	}
	return nil, nil
}

func (m *MockTradovateClient) GetAccounts() ([]models.Account, error) {
	if m.getAccountsFunc != nil {
		return m.getAccountsFunc()
	}
	return nil, nil
}

func (m *MockTradovateClient) GetRiskLimits(accountID int) (*models.RiskLimit, error) {
	if m.getRiskLimitsFunc != nil {
		return m.getRiskLimitsFunc(accountID)
	}
	return nil, nil
}

func (m *MockTradovateClient) PlaceOrder(order models.Order) (*models.Order, error) {
	if m.placeOrderFunc != nil {
		return m.placeOrderFunc(order)
	}
	return nil, nil
}

func (m *MockTradovateClient) CancelOrder(orderID int) error {
	if m.cancelOrderFunc != nil {
		return m.cancelOrderFunc(orderID)
	}
	return nil
}

func (m *MockTradovateClient) GetFills(orderID int) ([]models.Fill, error) {
	if m.getFillsFunc != nil {
		return m.getFillsFunc(orderID)
	}
	return nil, nil
}

func (m *MockTradovateClient) GetPositions() ([]models.Position, error) {
	if m.getPositionsFunc != nil {
		return m.getPositionsFunc()
	}
	return nil, nil
}

func (m *MockTradovateClient) GetContracts() ([]models.Contract, error) {
	if m.getContractsFunc != nil {
		return m.getContractsFunc()
	}
	return nil, nil
}

func (m *MockTradovateClient) GetMarketData(contractID int) (*models.MarketData, error) {
	if m.getMarketDataFunc != nil {
		return m.getMarketDataFunc(contractID)
	}
	return nil, nil
}

func (m *MockTradovateClient) GetHistoricalData(contractID int, startTime, endTime time.Time, interval string) ([]models.HistoricalData, error) {
	if m.getHistoricalDataFunc != nil {
		return m.getHistoricalDataFunc(contractID, startTime, endTime, interval)
	}
	return []models.HistoricalData{
		{
			ContractID: contractID,
			Timestamp:  startTime.Unix(),
			Open:       100.0,
			High:       101.0,
			Low:        99.0,
			Close:      100.5,
			Volume:     1000,
		},
	}, nil
}

func TestHandleAuthenticate(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() (*client.AuthResponse, error)
		wantErr bool
	}{
		{
			name: "Successful authentication",
			mockFn: func() (*client.AuthResponse, error) {
				return &client.AuthResponse{
					AccessToken:    "test-token",
					MdAccessToken:  "test-md-token",
					ExpirationTime: "2024-12-31T23:59:59Z",
					UserID:         12345,
					Name:           "Test User",
				}, nil
			},
			wantErr: false,
		},
		{
			name: "Authentication failure",
			mockFn: func() (*client.AuthResponse, error) {
				return nil, errors.New("invalid credentials")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockTradovateClient{
				authenticateFunc: tt.mockFn,
			}
			handlers := NewHandlers(mockClient)
			authHandler := handlers["authenticate"]

			result, err := authHandler.Handler(nil)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestHandleSetRiskLimits(t *testing.T) {
	errInvalidAccount := errors.New("invalid account ID")

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		mockErr error
		errMsg  string
	}{
		{
			name: "Valid risk limits",
			params: map[string]interface{}{
				"accountId":      float64(12345),
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(50.0),
			},
			wantErr: false,
			mockErr: nil,
		},
		{
			name: "Invalid account ID",
			params: map[string]interface{}{
				"accountId":      float64(-1),
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(50.0),
			},
			wantErr: true,
			mockErr: errInvalidAccount,
		},
		{
			name: "Missing account ID",
			params: map[string]interface{}{
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(50.0),
			},
			wantErr: true,
			errMsg:  "missing or invalid accountId",
		},
		{
			name: "Invalid day max loss",
			params: map[string]interface{}{
				"accountId":      float64(12345),
				"dayMaxLoss":     float64(-1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(50.0),
			},
			wantErr: true,
			errMsg:  "missing or invalid dayMaxLoss",
		},
		{
			name: "Invalid max drawdown",
			params: map[string]interface{}{
				"accountId":      float64(12345),
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(-500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(50.0),
			},
			wantErr: true,
			errMsg:  "missing or invalid maxDrawdown",
		},
		{
			name: "Invalid max position quantity",
			params: map[string]interface{}{
				"accountId":      float64(12345),
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(-10),
				"trailingStop":   float64(50.0),
			},
			wantErr: true,
			errMsg:  "missing or invalid maxPositionQty",
		},
		{
			name: "Invalid trailing stop",
			params: map[string]interface{}{
				"accountId":      float64(12345),
				"dayMaxLoss":     float64(1000.0),
				"maxDrawdown":    float64(500.0),
				"maxPositionQty": float64(10),
				"trailingStop":   float64(-50.0),
			},
			wantErr: true,
			errMsg:  "missing or invalid trailingStop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockTradovateClient{
				setRiskLimitsFunc: func(limits models.RiskLimit) error {
					return tt.mockErr
				},
			}
			handlers := NewHandlers(mockClient)
			setRiskLimitsHandler := handlers["setRiskLimits"]

			_, err := setRiskLimitsHandler.Handler(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandlePlaceOrder(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		mockFn  func(models.Order) (*models.Order, error)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid order",
			params: map[string]interface{}{
				"accountId":   float64(12345),
				"contractId":  float64(54321),
				"orderType":   "Limit",
				"price":       float64(100.50),
				"quantity":    float64(10),
				"timeInForce": "Day",
			},
			mockFn: func(order models.Order) (*models.Order, error) {
				order.ID = 67890
				return &order, nil
			},
			wantErr: false,
		},
		{
			name: "Missing required fields",
			params: map[string]interface{}{
				"accountId": float64(12345),
				// Missing other required fields
			},
			mockFn: func(order models.Order) (*models.Order, error) {
				return nil, errors.New("missing required fields")
			},
			wantErr: true,
			errMsg:  "missing required field: contractId",
		},
		{
			name: "Invalid field type",
			params: map[string]interface{}{
				"accountId":   "12345", // String instead of float64
				"contractId":  float64(54321),
				"orderType":   "Limit",
				"price":       float64(100.50),
				"quantity":    float64(10),
				"timeInForce": "Day",
			},
			mockFn: func(order models.Order) (*models.Order, error) {
				return nil, errors.New("invalid field type")
			},
			wantErr: true,
			errMsg:  "invalid type assertion for accountId",
		},
		{
			name: "Missing price for limit order",
			params: map[string]interface{}{
				"accountId":   float64(12345),
				"contractId":  float64(54321),
				"orderType":   "Limit",
				"quantity":    float64(10),
				"timeInForce": "Day",
			},
			mockFn: func(order models.Order) (*models.Order, error) {
				return nil, errors.New("price required for limit order")
			},
			wantErr: true,
			errMsg:  "price is required for Limit orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockTradovateClient{
				placeOrderFunc: tt.mockFn,
			}
			handlers := NewHandlers(mockClient)
			placeOrderHandler := handlers["placeOrder"]

			result, err := placeOrderHandler.Handler(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				order := result.(*models.Order)
				assert.Equal(t, 67890, order.ID)
			}
		})
	}
}

func TestHandleCancelOrder(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		mockFn  func(int) error
		wantErr bool
	}{
		{
			name: "Valid cancel",
			params: map[string]interface{}{
				"orderId": float64(67890),
			},
			mockFn: func(orderID int) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "Invalid order ID",
			params: map[string]interface{}{
				"orderId": float64(-1),
			},
			mockFn: func(orderID int) error {
				return errors.New("invalid order ID")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockTradovateClient{
				cancelOrderFunc: tt.mockFn,
			}
			handlers := NewHandlers(mockClient)
			cancelOrderHandler := handlers["cancelOrder"]

			result, err := cancelOrderHandler.Handler(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				success := result.(map[string]bool)["success"]
				assert.True(t, success)
			}
		})
	}
}

func TestHandleGetFills(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		mockFn  func(int) ([]models.Fill, error)
		wantErr bool
	}{
		{
			name: "Valid fills request",
			params: map[string]interface{}{
				"orderId": float64(67890),
			},
			mockFn: func(orderID int) ([]models.Fill, error) {
				return []models.Fill{
					{
						ID:        1,
						OrderID:   67890,
						Price:     100.50,
						Quantity:  5,
						Timestamp: time.Now().Unix(),
					},
				}, nil
			},
			wantErr: false,
		},
		{
			name: "Invalid order ID",
			params: map[string]interface{}{
				"orderId": float64(-1),
			},
			mockFn: func(orderID int) ([]models.Fill, error) {
				return nil, errors.New("invalid order ID")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockTradovateClient{
				getFillsFunc: tt.mockFn,
			}
			handlers := NewHandlers(mockClient)
			getFillsHandler := handlers["getFills"]

			result, err := getFillsHandler.Handler(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				fills := result.([]models.Fill)
				assert.Len(t, fills, 1)
				assert.Equal(t, 67890, fills[0].OrderID)
			}
		})
	}
}

func TestNewHandlers(t *testing.T) {
	mockClient := &MockTradovateClient{}
	handlers := NewHandlers(mockClient)

	// Test all handler registrations
	expectedHandlers := []string{
		"authenticate",
		"getAccounts",
		"getPositions",
		"placeOrder",
		"cancelOrder",
		"getFills",
		"getContracts",
		"getMarketData",
		"getHistoricalData",
		"setRiskLimits",
		"getRiskLimits",
	}

	for _, name := range expectedHandlers {
		t.Run(name, func(t *testing.T) {
			handler, exists := handlers[name]
			assert.True(t, exists, "Handler %s should exist", name)
			assert.NotNil(t, handler.Description, "Handler %s should have a description", name)
			assert.NotNil(t, handler.Handler, "Handler %s should have a handler function", name)
		})
	}
}

func TestGetAccountsHandler(t *testing.T) {
	mockAccounts := []models.Account{
		{ID: 1, Name: "Test Account"},
	}

	mockClient := &MockTradovateClient{
		getAccountsFunc: func() ([]models.Account, error) {
			return mockAccounts, nil
		},
	}

	handlers := NewHandlers(mockClient)
	result, err := handlers["getAccounts"].Handler(nil)
	assert.NoError(t, err)
	assert.Equal(t, mockAccounts, result)
}

func TestGetPositionsHandler(t *testing.T) {
	mockPositions := []models.Position{
		{ID: 1, AccountID: 123},
	}

	mockClient := &MockTradovateClient{
		getPositionsFunc: func() ([]models.Position, error) {
			return mockPositions, nil
		},
	}

	handlers := NewHandlers(mockClient)
	result, err := handlers["getPositions"].Handler(nil)
	assert.NoError(t, err)
	assert.Equal(t, mockPositions, result)
}

func TestGetContractsHandler(t *testing.T) {
	mockContracts := []models.Contract{
		{ID: 1, Name: "Test Contract"},
	}

	mockClient := &MockTradovateClient{
		getContractsFunc: func() ([]models.Contract, error) {
			return mockContracts, nil
		},
	}

	handlers := NewHandlers(mockClient)
	result, err := handlers["getContracts"].Handler(nil)
	assert.NoError(t, err)
	assert.Equal(t, mockContracts, result)
}

func TestGetMarketDataHandler(t *testing.T) {
	mockMarketData := &models.MarketData{
		ContractID: 1,
		Bid:        100.0,
		Ask:        101.0,
	}

	mockClient := &MockTradovateClient{
		getMarketDataFunc: func(contractID int) (*models.MarketData, error) {
			assert.Equal(t, 1, contractID)
			return mockMarketData, nil
		},
	}

	handlers := NewHandlers(mockClient)
	result, err := handlers["getMarketData"].Handler(map[string]interface{}{
		"contractId": float64(1),
	})
	assert.NoError(t, err)
	assert.Equal(t, mockMarketData, result)
}

func TestGetHistoricalDataHandler(t *testing.T) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	handlers := NewHandlers(&MockTradovateClient{})
	result, err := handlers["getHistoricalData"].Handler(map[string]interface{}{
		"contractId": float64(1),
		"startTime":  startTime.Format(time.RFC3339),
		"endTime":    endTime.Format(time.RFC3339),
		"interval":   "1h",
	})

	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.NotNil(t, result)
	}
}

func TestGetRiskLimitsHandler(t *testing.T) {
	expectedLimits := &models.RiskLimit{
		AccountID:      1,
		DayMaxLoss:     1000.0,
		MaxDrawdown:    500.0,
		MaxPositionQty: 10,
		TrailingStop:   50.0,
	}

	mockClient := &MockTradovateClient{
		getRiskLimitsFunc: func(accountID int) (*models.RiskLimit, error) {
			assert.Equal(t, 1, accountID)
			return expectedLimits, nil
		},
	}

	handlers := NewHandlers(mockClient)
	result, err := handlers["getRiskLimits"].Handler(map[string]interface{}{
		"accountId": float64(1),
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedLimits, result)
}

func TestHandleGetMarketDataInvalidParams(t *testing.T) {
	mockClient := &MockTradovateClient{}
	handlers := NewHandlers(mockClient)

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Missing contract ID",
			params:  map[string]interface{}{},
			wantErr: true,
			errMsg:  "missing contractId",
		},
		{
			name: "Invalid contract ID type",
			params: map[string]interface{}{
				"contractId": "invalid",
			},
			wantErr: true,
			errMsg:  "invalid type assertion for contractId",
		},
		{
			name: "Negative contract ID",
			params: map[string]interface{}{
				"contractId": float64(-1),
			},
			wantErr: true,
			errMsg:  "invalid contractId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handlers["getMarketData"].Handler(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleGetHistoricalDataInvalidParams(t *testing.T) {
	mockClient := &MockTradovateClient{}
	handlers := NewHandlers(mockClient)

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Missing all parameters",
			params:  map[string]interface{}{},
			wantErr: true,
			errMsg:  "missing contractId",
		},
		{
			name: "Invalid contract ID type",
			params: map[string]interface{}{
				"contractId": "invalid",
				"startTime":  time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"endTime":    time.Now().Format(time.RFC3339),
				"interval":   "1h",
			},
			wantErr: true,
			errMsg:  "invalid type assertion for contractId",
		},
		{
			name: "Invalid start time",
			params: map[string]interface{}{
				"contractId": float64(1),
				"startTime":  "invalid",
				"endTime":    time.Now().Format(time.RFC3339),
				"interval":   "1h",
			},
			wantErr: true,
			errMsg:  "invalid start time",
		},
		{
			name: "Invalid end time",
			params: map[string]interface{}{
				"contractId": float64(1),
				"startTime":  time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"endTime":    "invalid",
				"interval":   "1h",
			},
			wantErr: true,
			errMsg:  "invalid end time",
		},
		{
			name: "Missing interval",
			params: map[string]interface{}{
				"contractId": float64(1),
				"startTime":  time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"endTime":    time.Now().Format(time.RFC3339),
			},
			wantErr: true,
			errMsg:  "missing interval",
		},
		{
			name: "End time before start time",
			params: map[string]interface{}{
				"contractId": float64(1),
				"startTime":  time.Now().Format(time.RFC3339),
				"endTime":    time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"interval":   "1h",
			},
			wantErr: true,
			errMsg:  "end time must be after start time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handlers["getHistoricalData"].Handler(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleGetRiskLimitsInvalidParams(t *testing.T) {
	mockClient := &MockTradovateClient{}
	handlers := NewHandlers(mockClient)

	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Missing account ID",
			params:  map[string]interface{}{},
			wantErr: true,
			errMsg:  "missing accountId",
		},
		{
			name: "Invalid account ID type",
			params: map[string]interface{}{
				"accountId": "invalid",
			},
			wantErr: true,
			errMsg:  "invalid type assertion for accountId",
		},
		{
			name: "Negative account ID",
			params: map[string]interface{}{
				"accountId": float64(-1),
			},
			wantErr: true,
			errMsg:  "invalid accountId",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handlers["getRiskLimits"].Handler(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
