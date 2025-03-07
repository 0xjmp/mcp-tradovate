package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/0xjmp/mcp-tradovate/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewTradovateClient(t *testing.T) {
	client := NewTradovateClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	assert.Equal(t, "https://live.tradovate.com/v1", client.baseURL)
}

func TestSetBaseURL(t *testing.T) {
	client := NewTradovateClient()
	client.SetBaseURL("http://test-url")
	assert.Equal(t, "http://test-url", client.baseURL)
}

func TestAuthenticate(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/auth/accessTokenRequest", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Decode request body
		var authReq AuthRequest
		err := json.NewDecoder(r.Body).Decode(&authReq)
		assert.NoError(t, err)

		// Verify credentials from environment
		assert.Equal(t, os.Getenv("TRADOVATE_USERNAME"), authReq.Name)

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

	// Create test client with server URL
	client := NewTradovateClient()
	client.SetBaseURL(server.URL)

	// Set test environment variables
	os.Setenv("TRADOVATE_USERNAME", "testuser")
	os.Setenv("TRADOVATE_PASSWORD", "testpass")
	os.Setenv("TRADOVATE_APP_ID", "testapp")
	os.Setenv("TRADOVATE_APP_VERSION", "1.0")
	os.Setenv("TRADOVATE_CID", "testcid")
	os.Setenv("TRADOVATE_SEC", "testsec")

	// Test authentication
	authResp, err := client.Authenticate()
	assert.NoError(t, err)
	assert.Equal(t, "test-token", authResp.AccessToken)
	assert.Equal(t, 12345, authResp.UserID)
	assert.Equal(t, "test-token", client.GetAccessToken())
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

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)

	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Equal(t, "authentication failed: Invalid credentials", err.Error())
}

func TestAuthenticateNetworkError(t *testing.T) {
	client := NewTradovateClient()
	client.SetBaseURL("http://invalid-url")

	_, err := client.Authenticate()
	assert.Error(t, err)
}

func TestGetAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/account/list", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		accounts := []models.Account{
			{
				ID:            1,
				Name:          "Test Account",
				AccountType:   "Demo",
				Active:        true,
				CashBalance:   1000.0,
				RealizedPnL:   100.0,
				UnrealizedPnL: -50.0,
			},
		}
		json.NewEncoder(w).Encode(accounts)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	accounts, err := client.GetAccounts()
	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, "Test Account", accounts[0].Name)
}

func TestSetRiskLimits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/account/setRiskLimits", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var limits models.RiskLimit
		err := json.NewDecoder(r.Body).Decode(&limits)
		assert.NoError(t, err)
		assert.Equal(t, 12345, limits.AccountID)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	limits := models.RiskLimit{
		AccountID:      12345,
		DayMaxLoss:     1000.0,
		MaxDrawdown:    500.0,
		MaxPositionQty: 10,
		TrailingStop:   50.0,
	}

	err := client.SetRiskLimits(limits)
	assert.NoError(t, err)
}

func TestGetRiskLimits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/account/riskLimits/12345", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		limits := models.RiskLimit{
			AccountID:      12345,
			DayMaxLoss:     1000.0,
			MaxDrawdown:    500.0,
			MaxPositionQty: 10,
			TrailingStop:   50.0,
		}
		json.NewEncoder(w).Encode(limits)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	limits, err := client.GetRiskLimits(12345)
	assert.NoError(t, err)
	assert.Equal(t, 12345, limits.AccountID)
	assert.Equal(t, 1000.0, limits.DayMaxLoss)
}

func TestPlaceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/order/placeOrder", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		var order models.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		assert.NoError(t, err)
		assert.Equal(t, 12345, order.AccountID)

		order.ID = 67890 // Add order ID in response
		json.NewEncoder(w).Encode(order)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	order := models.Order{
		AccountID:   12345,
		ContractID:  54321,
		OrderType:   "Limit",
		Price:       100.50,
		Quantity:    10,
		TimeInForce: "Day",
	}

	placedOrder, err := client.PlaceOrder(order)
	assert.NoError(t, err)
	assert.Equal(t, 67890, placedOrder.ID)
	assert.Equal(t, order.AccountID, placedOrder.AccountID)
}

func TestCancelOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/order/cancel/67890", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	err := client.CancelOrder(67890)
	assert.NoError(t, err)
}

func TestGetFills(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/fill/list/67890", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		fills := []models.Fill{
			{
				ID:        1,
				OrderID:   67890,
				Price:     100.50,
				Quantity:  5,
				Timestamp: time.Now().Unix(),
			},
		}
		json.NewEncoder(w).Encode(fills)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	fills, err := client.GetFills(67890)
	assert.NoError(t, err)
	assert.Len(t, fills, 1)
	assert.Equal(t, 67890, fills[0].OrderID)
}

func TestGetPositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/position/list", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		positions := []models.Position{
			{
				ID:           1,
				AccountID:    12345,
				ContractID:   54321,
				NetPos:       5,
				AvgPrice:     100.50,
				RealizedPL:   250.75,
				UnrealizedPL: -125.50,
			},
		}
		json.NewEncoder(w).Encode(positions)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	positions, err := client.GetPositions()
	assert.NoError(t, err)
	assert.Len(t, positions, 1)
	assert.Equal(t, 5, positions[0].NetPos)
}

func TestGetContracts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/contract/list", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		contracts := []models.Contract{
			{
				ID:           1,
				Name:         "ES Mar24",
				ContractType: "Future",
				Exchange:     "CME",
				Symbol:       "ESH4",
			},
		}
		json.NewEncoder(w).Encode(contracts)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	contracts, err := client.GetContracts()
	assert.NoError(t, err)
	assert.Len(t, contracts, 1)
	assert.Equal(t, "ES Mar24", contracts[0].Name)
}

func TestGetMarketData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/md/getQuote/54321", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		marketData := models.MarketData{
			ContractID: 54321,
			Bid:        100.25,
			Ask:        100.50,
			Last:       100.75,
			Volume:     1500,
			Timestamp:  time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(marketData)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	data, err := client.GetMarketData(54321)
	assert.NoError(t, err)
	assert.Equal(t, 54321, data.ContractID)
	assert.Equal(t, 100.25, data.Bid)
}

func TestGetHistoricalData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/md/historical", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		data := []models.HistoricalData{
			{
				ContractID: 54321,
				Timestamp:  time.Now().Unix(),
				Open:       100.25,
				High:       101.50,
				Low:        99.75,
				Close:      100.50,
				Volume:     1500,
			},
		}
		json.NewEncoder(w).Encode(data)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	data, err := client.GetHistoricalData(54321, startTime, endTime, "1h")
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, 54321, data[0].ContractID)
}

func TestDoRequestError(t *testing.T) {
	client := NewTradovateClient()
	client.SetBaseURL("http://invalid-url")

	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		wantErr bool
	}{
		{
			name:    "Invalid URL",
			method:  "GET",
			path:    "/test",
			body:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid request method",
			method:  "\n",
			path:    "/test",
			body:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid JSON body",
			method:  "POST",
			path:    "/test",
			body:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.doRequest(tt.method, tt.path, tt.body)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClientErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"errorText": "Internal server error",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	tests := []struct {
		name    string
		test    func() error
		wantErr bool
	}{
		{
			name: "GetAccounts error",
			test: func() error {
				_, err := client.GetAccounts()
				return err
			},
			wantErr: true,
		},
		{
			name: "GetRiskLimits error",
			test: func() error {
				_, err := client.GetRiskLimits(1)
				return err
			},
			wantErr: true,
		},
		{
			name: "SetRiskLimits error",
			test: func() error {
				return client.SetRiskLimits(models.RiskLimit{})
			},
			wantErr: true,
		},
		{
			name: "PlaceOrder error",
			test: func() error {
				_, err := client.PlaceOrder(models.Order{})
				return err
			},
			wantErr: true,
		},
		{
			name: "CancelOrder error",
			test: func() error {
				return client.CancelOrder(1)
			},
			wantErr: true,
		},
		{
			name: "GetFills error",
			test: func() error {
				_, err := client.GetFills(1)
				return err
			},
			wantErr: true,
		},
		{
			name: "GetPositions error",
			test: func() error {
				_, err := client.GetPositions()
				return err
			},
			wantErr: true,
		},
		{
			name: "GetContracts error",
			test: func() error {
				_, err := client.GetContracts()
				return err
			},
			wantErr: true,
		},
		{
			name: "GetMarketData error",
			test: func() error {
				_, err := client.GetMarketData(1)
				return err
			},
			wantErr: true,
		},
		{
			name: "GetHistoricalData error",
			test: func() error {
				_, err := client.GetHistoricalData(1, time.Now(), time.Now(), "1h")
				return err
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.test()
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), "status 500")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthenticateInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/auth/accessTokenRequest", r.URL.Path)

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"errorText": "Invalid credentials",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)

	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid credentials")
}

func TestAuthenticateInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)

	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
}

func TestNetworkTimeoutError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.httpClient.Timeout = 1 * time.Second

	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestInvalidParameterValidation(t *testing.T) {
	client := NewTradovateClient()

	// Test invalid order
	order := models.Order{
		AccountID:  -1,
		OrderType:  "InvalidType",
		Quantity:   -10,
		ContractID: -1,
	}
	_, err := client.PlaceOrder(order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")

	// Test invalid risk limits
	limits := models.RiskLimit{
		AccountID:      -1,
		DayMaxLoss:     -1000,
		MaxDrawdown:    -500,
		MaxPositionQty: -10,
		TrailingStop:   -50,
	}
	err = client.SetRiskLimits(limits)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")

	// Test invalid order ID
	err = client.CancelOrder(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")

	// Test invalid contract ID
	_, err = client.GetMarketData(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")

	// Test invalid historical data parameters
	_, err = client.GetHistoricalData(-1, time.Now(), time.Now().Add(-24*time.Hour), "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 400")
}

func TestInvalidResponseHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)

	// Test various endpoints with invalid responses
	_, err := client.Authenticate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")

	_, err = client.GetAccounts()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")

	_, err = client.GetPositions()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")

	_, err = client.GetContracts()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")

	_, err = client.GetMarketData(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")

	_, err = client.GetHistoricalData(1, time.Now().Add(-24*time.Hour), time.Now(), "1h")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestGetAccountsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	_, err := client.GetAccounts()
	assert.Error(t, err)
}

func TestSetRiskLimitsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid risk limits",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	limits := models.RiskLimit{
		AccountID: 12345,
	}

	err := client.SetRiskLimits(limits)
	assert.Error(t, err)
}

func TestGetRiskLimitsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Risk limits not found",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	_, err := client.GetRiskLimits(12345)
	assert.Error(t, err)
}

func TestPlaceOrderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid order parameters",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	order := models.Order{
		AccountID: 12345,
	}

	_, err := client.PlaceOrder(order)
	assert.Error(t, err)
}

func TestCancelOrderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order not found",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	err := client.CancelOrder(67890)
	assert.Error(t, err)
}

func TestGetFillsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
	}))
	defer server.Close()

	client := NewTradovateClient()
	client.SetBaseURL(server.URL)
	client.accessToken = "test-token"

	_, err := client.GetFills(67890)
	assert.Error(t, err)
}

func TestNetworkErrors(t *testing.T) {
	client := NewTradovateClient()
	client.SetBaseURL("http://invalid-url")
	client.accessToken = "test-token"

	// Test GetAccounts network error
	_, err := client.GetAccounts()
	assert.Error(t, err)

	// Test SetRiskLimits network error
	err = client.SetRiskLimits(models.RiskLimit{AccountID: 12345})
	assert.Error(t, err)

	// Test GetRiskLimits network error
	_, err = client.GetRiskLimits(12345)
	assert.Error(t, err)

	// Test PlaceOrder network error
	_, err = client.PlaceOrder(models.Order{AccountID: 12345})
	assert.Error(t, err)

	// Test CancelOrder network error
	err = client.CancelOrder(67890)
	assert.Error(t, err)

	// Test GetFills network error
	_, err = client.GetFills(67890)
	assert.Error(t, err)
}
