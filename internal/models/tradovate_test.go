package models

import (
	"encoding/json"
	"testing"
)

func TestAccountMarshaling(t *testing.T) {
	account := Account{
		ID:            12345,
		Name:          "Test Account",
		AccountType:   "Demo",
		Active:        true,
		CashBalance:   1000.50,
		RealizedPnL:   100.25,
		UnrealizedPnL: -50.75,
	}

	// Test marshaling
	data, err := json.Marshal(account)
	if err != nil {
		t.Errorf("Failed to marshal Account: %v", err)
	}

	// Test unmarshaling
	var decoded Account
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal Account: %v", err)
	}

	// Verify fields
	if decoded.ID != account.ID {
		t.Errorf("Expected ID %d, got %d", account.ID, decoded.ID)
	}
	if decoded.Name != account.Name {
		t.Errorf("Expected Name %s, got %s", account.Name, decoded.Name)
	}
	if decoded.CashBalance != account.CashBalance {
		t.Errorf("Expected CashBalance %f, got %f", account.CashBalance, decoded.CashBalance)
	}
}

func TestOrderMarshaling(t *testing.T) {
	order := Order{
		ID:          67890,
		AccountID:   12345,
		ContractID:  54321,
		OrderType:   "Limit",
		Status:      "Working",
		Side:        "Buy",
		Price:       100.50,
		StopPrice:   0,
		Quantity:    10,
		FilledQty:   5,
		TimeInForce: "Day",
	}

	// Test marshaling
	data, err := json.Marshal(order)
	if err != nil {
		t.Errorf("Failed to marshal Order: %v", err)
	}

	// Test unmarshaling
	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal Order: %v", err)
	}

	// Verify fields
	if decoded.ID != order.ID {
		t.Errorf("Expected ID %d, got %d", order.ID, decoded.ID)
	}
	if decoded.OrderType != order.OrderType {
		t.Errorf("Expected OrderType %s, got %s", order.OrderType, decoded.OrderType)
	}
	if decoded.Quantity != order.Quantity {
		t.Errorf("Expected Quantity %d, got %d", order.Quantity, decoded.Quantity)
	}
}

func TestPositionMarshaling(t *testing.T) {
	position := Position{
		ID:           13579,
		AccountID:    12345,
		ContractID:   54321,
		NetPos:       5,
		AvgPrice:     100.50,
		RealizedPL:   250.75,
		UnrealizedPL: -125.50,
	}

	// Test marshaling
	data, err := json.Marshal(position)
	if err != nil {
		t.Errorf("Failed to marshal Position: %v", err)
	}

	// Test unmarshaling
	var decoded Position
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal Position: %v", err)
	}

	// Verify fields
	if decoded.ID != position.ID {
		t.Errorf("Expected ID %d, got %d", position.ID, decoded.ID)
	}
	if decoded.NetPos != position.NetPos {
		t.Errorf("Expected NetPos %d, got %d", position.NetPos, decoded.NetPos)
	}
	if decoded.AvgPrice != position.AvgPrice {
		t.Errorf("Expected AvgPrice %f, got %f", position.AvgPrice, decoded.AvgPrice)
	}
}

func TestContractMarshaling(t *testing.T) {
	contract := Contract{
		ID:           24680,
		Name:         "ES Mar24",
		ContractType: "Future",
		Exchange:     "CME",
		Symbol:       "ESH4",
	}

	// Test marshaling
	data, err := json.Marshal(contract)
	if err != nil {
		t.Errorf("Failed to marshal Contract: %v", err)
	}

	// Test unmarshaling
	var decoded Contract
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal Contract: %v", err)
	}

	// Verify fields
	if decoded.ID != contract.ID {
		t.Errorf("Expected ID %d, got %d", contract.ID, decoded.ID)
	}
	if decoded.Name != contract.Name {
		t.Errorf("Expected Name %s, got %s", contract.Name, decoded.Name)
	}
	if decoded.Symbol != contract.Symbol {
		t.Errorf("Expected Symbol %s, got %s", contract.Symbol, decoded.Symbol)
	}
}

func TestMarketDataMarshaling(t *testing.T) {
	marketData := MarketData{
		ContractID: 54321,
		Bid:        100.25,
		Ask:        100.50,
		Last:       100.75,
		Volume:     1500,
		Timestamp:  1709876543,
	}

	// Test marshaling
	data, err := json.Marshal(marketData)
	if err != nil {
		t.Errorf("Failed to marshal MarketData: %v", err)
	}

	// Test unmarshaling
	var decoded MarketData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal MarketData: %v", err)
	}

	// Verify fields
	if decoded.ContractID != marketData.ContractID {
		t.Errorf("Expected ContractID %d, got %d", marketData.ContractID, decoded.ContractID)
	}
	if decoded.Bid != marketData.Bid {
		t.Errorf("Expected Bid %f, got %f", marketData.Bid, decoded.Bid)
	}
	if decoded.Volume != marketData.Volume {
		t.Errorf("Expected Volume %d, got %d", marketData.Volume, decoded.Volume)
	}
}

func TestHistoricalDataMarshaling(t *testing.T) {
	historicalData := HistoricalData{
		ContractID: 54321,
		Timestamp:  1709876543,
		Open:       100.25,
		High:       101.50,
		Low:        99.75,
		Close:      100.50,
		Volume:     1500,
	}

	// Test marshaling
	data, err := json.Marshal(historicalData)
	if err != nil {
		t.Errorf("Failed to marshal HistoricalData: %v", err)
	}

	// Test unmarshaling
	var decoded HistoricalData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("Failed to unmarshal HistoricalData: %v", err)
	}

	// Verify fields
	if decoded.ContractID != historicalData.ContractID {
		t.Errorf("Expected ContractID %d, got %d", historicalData.ContractID, decoded.ContractID)
	}
	if decoded.Open != historicalData.Open {
		t.Errorf("Expected Open %f, got %f", historicalData.Open, decoded.Open)
	}
	if decoded.Volume != historicalData.Volume {
		t.Errorf("Expected Volume %d, got %d", historicalData.Volume, decoded.Volume)
	}
}
