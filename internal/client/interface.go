package client

import (
	"time"

	"github.com/0jxmp/tradovate-mcp/internal/models"
)

// TradovateClient defines the interface for interacting with the Tradovate API
type TradovateClient interface {
	// Authentication
	Authenticate(creds models.Credentials) (*models.AccessToken, error)

	// Account Operations
	GetAccounts() ([]models.Account, error)
	GetRiskLimits(accountID int) (*models.RiskLimit, error)
	SetRiskLimits(limits models.RiskLimit) error

	// Trading Operations
	PlaceOrder(order models.Order) (*models.Order, error)
	CancelOrder(orderID int) error
	GetFills(orderID int) ([]models.Fill, error)
	GetPositions() ([]models.Position, error)

	// Market Data Operations
	GetContracts() ([]models.Contract, error)
	GetMarketData(contractID int) (*models.MarketData, error)
	GetHistoricalData(contractID int, startTime, endTime time.Time, interval string) ([]models.HistoricalData, error)
}
