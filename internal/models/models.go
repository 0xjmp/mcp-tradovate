package models

import "time"

// Credentials represents authentication credentials
type Credentials struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	AppID      string `json:"appId"`
	AppVersion string `json:"appVersion"`
	CID        string `json:"cid"`
	SEC        string `json:"sec"`
}

// AccessToken represents the authentication response
type AccessToken struct {
	AccessToken    string `json:"accessToken"`
	ExpirationTime int64  `json:"expirationTime"`
	UserID         int    `json:"userId"`
	Name           string `json:"name"`
}

// Position represents a trading position
type Position struct {
	ContractID  int     `json:"contractId"`
	AccountID   int     `json:"accountId"`
	Timestamp   int64   `json:"timestamp"`
	TradePrice  float64 `json:"tradePrice"`
	Position    int     `json:"position"`
	ProfitLoss  float64 `json:"profitLoss"`
	OpenPrice   float64 `json:"openPrice"`
	MarketValue float64 `json:"marketValue"`
}

// Order represents a trading order
type Order struct {
	ID             int     `json:"id"`
	AccountID      int     `json:"accountId"`
	ContractID     int     `json:"contractId"`
	OrderType      string  `json:"orderType"`
	Price          float64 `json:"price"`
	Quantity       int     `json:"quantity"`
	Status         string  `json:"status"`
	FilledQty      int     `json:"filledQty"`
	TimeInForce    string  `json:"timeInForce"`
	ExpirationTime int64   `json:"expirationTime,omitempty"`
}

// Contract represents a tradable contract
type Contract struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Exchange    string `json:"exchange"`
	Symbol      string `json:"symbol"`
}

// Account represents a trading account
type Account struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	UserID        int     `json:"userId"`
	AccountType   string  `json:"accountType"`
	Active        bool    `json:"active"`
	CashBalance   float64 `json:"cashBalance"`
	RealizedPnL   float64 `json:"realizedPnL"`
	UnrealizedPnL float64 `json:"unrealizedPnL"`
	MarginUsed    float64 `json:"marginUsed"`
}

// MarketData represents real-time market data
type MarketData struct {
	ContractID int       `json:"contractId"`
	Timestamp  time.Time `json:"timestamp"`
	Bid        float64   `json:"bid"`
	Ask        float64   `json:"ask"`
	Last       float64   `json:"last"`
	Volume     int       `json:"volume"`
}

// Fill represents an order fill
type Fill struct {
	ID         int       `json:"id"`
	OrderID    int       `json:"orderId"`
	Timestamp  time.Time `json:"timestamp"`
	Price      float64   `json:"price"`
	Quantity   int       `json:"quantity"`
	Commission float64   `json:"commission"`
}

// HistoricalData represents historical price data
type HistoricalData struct {
	ContractID int       `json:"contractId"`
	Timestamp  time.Time `json:"timestamp"`
	Open       float64   `json:"open"`
	High       float64   `json:"high"`
	Low        float64   `json:"low"`
	Close      float64   `json:"close"`
	Volume     int       `json:"volume"`
}

// RiskLimit represents account risk management settings
type RiskLimit struct {
	AccountID     int     `json:"accountId"`
	MaxPositions  int     `json:"maxPositions"`
	MaxLoss       float64 `json:"maxLoss"`
	DailyMaxLoss  float64 `json:"dailyMaxLoss"`
	MarginPercent float64 `json:"marginPercent"`
}
