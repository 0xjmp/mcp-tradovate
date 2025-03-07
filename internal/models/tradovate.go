package models

// Account represents a Tradovate trading account
type Account struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	AccountType   string  `json:"accountType"`
	Active        bool    `json:"active"`
	CashBalance   float64 `json:"cashBalance"`
	RealizedPnL   float64 `json:"realizedPnL"`
	UnrealizedPnL float64 `json:"unrealizedPnL"`
}

// Order represents a trading order
type Order struct {
	ID          int     `json:"id"`
	AccountID   int     `json:"accountId"`
	ContractID  int     `json:"contractId"`
	OrderType   string  `json:"orderType"`
	Status      string  `json:"status"`
	Side        string  `json:"side"`
	Price       float64 `json:"price,omitempty"`
	StopPrice   float64 `json:"stopPrice,omitempty"`
	Quantity    int     `json:"qty"`
	FilledQty   int     `json:"filledQty"`
	TimeInForce string  `json:"timeInForce"`
}

// Position represents a trading position
type Position struct {
	ID           int     `json:"id"`
	AccountID    int     `json:"accountId"`
	ContractID   int     `json:"contractId"`
	NetPos       int     `json:"netPos"`
	AvgPrice     float64 `json:"avgPrice"`
	RealizedPL   float64 `json:"realizedPL"`
	UnrealizedPL float64 `json:"unrealizedPL"`
}

// Contract represents a tradable contract
type Contract struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ContractType string `json:"contractType"`
	Exchange     string `json:"exchange"`
	Symbol       string `json:"symbol"`
}

// MarketData represents real-time market data
type MarketData struct {
	ContractID int     `json:"contractId"`
	Bid        float64 `json:"bid"`
	Ask        float64 `json:"ask"`
	Last       float64 `json:"last"`
	Volume     int     `json:"volume"`
	Timestamp  int64   `json:"timestamp"`
}

// Fill represents a trade fill
type Fill struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"orderId"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"qty"`
	Timestamp int64   `json:"timestamp"`
}

// RiskLimit represents account risk limits
type RiskLimit struct {
	AccountID      int     `json:"accountId"`
	DayMaxLoss     float64 `json:"dayMaxLoss"`
	MaxDrawdown    float64 `json:"maxDrawdown"`
	MaxPositionQty int     `json:"maxPositionQty"`
	TrailingStop   float64 `json:"trailingStop"`
}

// HistoricalData represents historical price data
type HistoricalData struct {
	ContractID int     `json:"contractId"`
	Timestamp  int64   `json:"timestamp"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     int     `json:"volume"`
}
