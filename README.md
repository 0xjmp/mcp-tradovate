# Tradovate MCP Server

[![smithery](https://smithery.ai/badge.svg)](https://smithery.ai)

An MCP server implementation that provides tools for interacting with the Tradovate trading platform.

## Features

* Authentication with Tradovate API
* Account management and risk controls
* Real-time market data access
* Order placement and management
* Position tracking
* Historical data retrieval

## Tool

### tradovate

Facilitates interaction with the Tradovate trading platform.

**Methods:**

* `authenticate`: Authenticate with Tradovate using provided credentials
* `getAccounts`: Retrieve trading accounts
* `getRiskLimits`: Get account risk limits
* `setRiskLimits`: Update account risk limits
* `placeOrder`: Place a new trading order
* `cancelOrder`: Cancel an existing order
* `getFills`: Get order fill information
* `getPositions`: Retrieve current positions
* `getContracts`: Get available trading contracts
* `getMarketData`: Get real-time market data
* `getHistoricalData`: Retrieve historical price data

## Configuration

The server requires the following environment variables:

```yaml
TRADOVATE_USERNAME: Tradovate account username
TRADOVATE_PASSWORD: Tradovate account password
TRADOVATE_APP_ID: Application ID from Tradovate
TRADOVATE_APP_VERSION: Application version (defaults to "1.0")
TRADOVATE_CID: Client ID from Tradovate
TRADOVATE_SEC: Client Secret from Tradovate
```

### Usage with Smithery

Add this to your `smithery.yaml`:

```yaml
startCommand:
  type: stdio
  configSchema:
    type: object
    required:
      - TRADOVATE_USERNAME
      - TRADOVATE_PASSWORD
      - TRADOVATE_APP_ID
      - TRADOVATE_CID
      - TRADOVATE_SEC
    properties:
      TRADOVATE_USERNAME:
        type: string
        description: Tradovate username
      TRADOVATE_PASSWORD:
        type: string
        description: Tradovate password
        secret: true
      TRADOVATE_APP_ID:
        type: string
        description: Tradovate Application ID
      TRADOVATE_APP_VERSION:
        type: string
        description: Tradovate Application Version
        default: "1.0"
      TRADOVATE_CID:
        type: string
        description: Tradovate Client ID
      TRADOVATE_SEC:
        type: string
        description: Tradovate Client Secret
        secret: true
```

## Building

Docker:
```bash
docker build -t mcp/tradovate .
```

Local:
```bash
go build -o mcp-tradovate ./cmd/mcp-tradovate
```

## Testing

Run the test suite:
```bash
go test ./...
```

## License

This MCP server is licensed under the MIT License. See the LICENSE file in the project repository for full license text. 