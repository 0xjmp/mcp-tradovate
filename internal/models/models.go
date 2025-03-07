// Package models defines the data structures used throughout the MCP Tradovate server.
// It provides type-safe representations of accounts, orders, positions, and market data
// that are used for communication with the Tradovate API.
package models

// Account represents a trading account in Tradovate.
type Account struct {
	ID            int     `json:"id"`            // Unique identifier for the account
	Name          string  `json:"name"`          // Account name
	AccountType   string  `json:"accountType"`   // Type of account (e.g., "Demo", "Live")
	Active        bool    `json:"active"`        // Whether the account is active
	CashBalance   float64 `json:"cashBalance"`   // Current cash balance
	RealizedPnL   float64 `json:"realizedPnL"`   // Realized profit and loss
	UnrealizedPnL float64 `json:"unrealizedPnL"` // Unrealized profit and loss
}

// Order represents a trading order in Tradovate.
type Order struct {
	ID           int     `json:"id,omitempty"`        // Unique identifier for the order
	AccountID    int     `json:"accountId"`           // Account that placed the order
	ContractID   int     `json:"contractId"`          // Contract being traded
	OrderType    string  `json:"orderType"`           // Type of order (Market, Limit, etc.)
	Side         string  `json:"side"`                // Order side (Buy, Sell)
	Price        float64 `json:"price"`               // Order price (required for Limit orders)
	StopPrice    float64 `json:"stopPrice,omitempty"` // Stop price for stop orders
	Quantity     int     `json:"quantity"`            // Number of contracts
	TimeInForce  string  `json:"timeInForce"`         // Time in force (Day, GTC, IOC, etc.)
	Status       string  `json:"status"`              // Current order status
	FilledQty    int     `json:"filledQty"`           // Number of contracts filled
	AveragePrice float64 `json:"averagePrice"`        // Average fill price
	CreatedAt    int64   `json:"createdAt"`           // Order creation timestamp
	UpdatedAt    int64   `json:"updatedAt"`           // Last update timestamp
}

// Fill represents an order fill in Tradovate.
type Fill struct {
	ID        int     `json:"id"`        // Unique identifier for the fill
	OrderID   int     `json:"orderId"`   // Order that was filled
	Price     float64 `json:"price"`     // Fill price
	Quantity  int     `json:"quantity"`  // Fill quantity
	Timestamp int64   `json:"timestamp"` // Fill timestamp
}

// Position represents a trading position in Tradovate.
type Position struct {
	ID           int     `json:"id"`           // Unique identifier for the position
	AccountID    int     `json:"accountId"`    // Account holding the position
	ContractID   int     `json:"contractId"`   // Contract being held
	NetPos       int     `json:"netPos"`       // Net position size
	AvgPrice     float64 `json:"avgPrice"`     // Average entry price
	RealizedPL   float64 `json:"realizedPL"`   // Realized profit/loss
	UnrealizedPL float64 `json:"unrealizedPL"` // Unrealized profit/loss
}

// Contract represents a tradable contract in Tradovate.
type Contract struct {
	ID           int    `json:"id"`           // Unique identifier for the contract
	Name         string `json:"name"`         // Contract name
	ContractType string `json:"contractType"` // Type of contract (Future, Option, etc.)
	Exchange     string `json:"exchange"`     // Exchange where contract is traded
	Symbol       string `json:"symbol"`       // Trading symbol
}

// MarketData represents real-time market data for a contract.
type MarketData struct {
	ContractID int     `json:"contractId"` // Contract this data is for
	Bid        float64 `json:"bid"`        // Best bid price
	Ask        float64 `json:"ask"`        // Best ask price
	Last       float64 `json:"last"`       // Last trade price
	Volume     int     `json:"volume"`     // Trading volume
	Timestamp  int64   `json:"timestamp"`  // Data timestamp
}

// HistoricalData represents historical price data for a contract.
type HistoricalData struct {
	ContractID int     `json:"contractId"` // Contract this data is for
	Timestamp  int64   `json:"timestamp"`  // Bar timestamp
	Open       float64 `json:"open"`       // Opening price
	High       float64 `json:"high"`       // Highest price
	Low        float64 `json:"low"`        // Lowest price
	Close      float64 `json:"close"`      // Closing price
	Volume     int     `json:"volume"`     // Trading volume
}

// RiskLimit represents risk management limits for an account.
type RiskLimit struct {
	AccountID      int     `json:"accountId"`      // Account these limits apply to
	DayMaxLoss     float64 `json:"dayMaxLoss"`     // Maximum loss allowed per day
	MaxDrawdown    float64 `json:"maxDrawdown"`    // Maximum drawdown allowed
	MaxPositionQty int     `json:"maxPositionQty"` // Maximum position size allowed
	TrailingStop   float64 `json:"trailingStop"`   // Trailing stop percentage
}
