# Tradovate MCP Plugin

This plugin integrates Tradovate's trading API with MCP, providing a comprehensive set of tools for trading automation and management. The plugin can be used with Claude Desktop app, Cursor, or through Smithery.ai.

## Project Structure

```
.
├── cmd/
│   └── tradovate-mcp/     # Main application entry point
├── internal/
│   ├── client/            # API client implementation
│   ├── handlers/          # MCP tool handlers
│   └── models/            # Domain models
├── pkg/                   # Public packages (if needed)
├── LICENSE
├── README.md
└── go.mod
```

## Architecture

The plugin follows clean architecture principles with clear separation of concerns:

- **Domain Models** (`internal/models`): Core business entities and data structures
- **Client Layer** (`internal/client`): Handles API communication with Tradovate
- **Handler Layer** (`internal/handlers`): Implements MCP tool handlers
- **Entry Point** (`cmd/tradovate-mcp`): Initializes and wires components together

## Features

### Authentication
- `authenticate`: Authenticate with Tradovate API using credentials

### Account Management
- `getAccounts`: Get all accounts for the authenticated user
- `getRiskLimits`: Get current risk management limits for an account
- `setRiskLimits`: Set risk management limits for an account

### Trading
- `placeOrder`: Place a new order
- `cancelOrder`: Cancel an existing order
- `getFills`: Get fills for a specific order
- `getPositions`: Get current positions

### Market Data
- `getContracts`: Get available contracts
- `getMarketData`: Get real-time market data for a contract
- `getHistoricalData`: Get historical price data for a contract

## Installation

### Option 1: Using Smithery.ai (Recommended)

#### For Cursor
1. Open Cursor
2. Go to Settings > Extensions
3. Click "Add Extension"
4. Search for "tradovate-mcp"
5. Click "Install"
6. Configure credentials in the extension settings

#### For Claude Desktop
1. Open Claude Desktop
2. Go to Settings > Plugins
3. Click "Browse Plugins"
4. Search for "tradovate-mcp"
5. Click "Install"
6. Configure credentials in the plugin settings

### Option 2: Manual Installation

#### Prerequisites
- Go 1.21 or later
- Git
- Make sure you have the correct plugin directory:
  ```bash
  # For Claude Desktop
  # macOS:   ~/Library/Application Support/Claude/plugins
  # Linux:   ~/.config/claude/plugins
  # Windows: %APPDATA%\Claude\plugins
  
  # For Cursor
  # macOS:   ~/Library/Application Support/Cursor/extensions
  # Linux:   ~/.config/cursor/extensions
  # Windows: %APPDATA%\Cursor\extensions
  ```

#### For Claude Desktop

1. Clone the repository to the plugins directory:
   ```bash
   # On macOS
   git clone https://github.com/0jxmp/tradovate-mcp.git ~/Library/Application\ Support/Claude/plugins/tradovate-mcp
   # On Linux
   git clone https://github.com/0jxmp/tradovate-mcp.git ~/.config/claude/plugins/tradovate-mcp
   # On Windows
   git clone https://github.com/0jxmp/tradovate-mcp.git %APPDATA%\Claude\plugins\tradovate-mcp
   ```

2. Install dependencies and build:
   ```bash
   cd ~/Library/Application\ Support/Claude/plugins/tradovate-mcp  # Adjust path based on your OS
   go mod download
   go build -o tradovate-mcp ./cmd/tradovate-mcp
   ```

3. Create the plugin configuration:
   ```bash
   cat > claude-plugin.json << EOF
   {
     "name": "tradovate-mcp",
     "version": "1.0.0",
     "description": "Tradovate API integration for Claude",
     "type": "mcp-plugin",
     "main": "tradovate-mcp",
     "autostart": true
   }
   EOF
   ```

4. Create the environment file:
   ```bash
   cat > .env << EOF
   TRADOVATE_USERNAME="your_username"
   TRADOVATE_PASSWORD="your_password"
   TRADOVATE_APP_ID="your_app_id"
   TRADOVATE_APP_VERSION="1.0"
   TRADOVATE_CID="your_client_id"
   TRADOVATE_SEC="your_client_secret"
   PORT="8080"
   EOF
   chmod 600 .env
   ```

5. Restart Claude Desktop

#### For Cursor

1. Clone the repository to the extensions directory:
   ```bash
   # On macOS
   git clone https://github.com/0jxmp/tradovate-mcp.git ~/Library/Application\ Support/Cursor/extensions/tradovate-mcp
   # On Linux
   git clone https://github.com/0jxmp/tradovate-mcp.git ~/.config/cursor/extensions/tradovate-mcp
   # On Windows
   git clone https://github.com/0jxmp/tradovate-mcp.git %APPDATA%\Cursor\extensions\tradovate-mcp
   ```

2. Install dependencies and build:
   ```bash
   cd ~/Library/Application\ Support/Cursor/extensions/tradovate-mcp  # Adjust path based on your OS
   go mod download
   go build -o tradovate-mcp ./cmd/tradovate-mcp
   ```

3. Create the extension configuration:
   ```bash
   cat > cursor-extension.json << EOF
   {
     "name": "tradovate-mcp",
     "version": "1.0.0",
     "description": "Tradovate API integration for Cursor",
     "type": "mcp-plugin",
     "main": "tradovate-mcp",
     "autostart": true,
     "capabilities": ["trading", "market_data"]
   }
   EOF
   ```

4. Create the environment file:
   ```bash
   cat > .env << EOF
   TRADOVATE_USERNAME="your_username"
   TRADOVATE_PASSWORD="your_password"
   TRADOVATE_APP_ID="your_app_id"
   TRADOVATE_APP_VERSION="1.0"
   TRADOVATE_CID="your_client_id"
   TRADOVATE_SEC="your_client_secret"
   PORT="8081"  # Use different port than Claude to avoid conflicts
   EOF
   chmod 600 .env
   ```

5. Restart Cursor

## Verifying Installation

### For Claude Desktop
1. Open Claude Desktop
2. The plugin should appear in the plugins list
3. Test the connection:
   ```go
   {
       "tool": "getContracts"
   }
   ```

### For Cursor
1. Open Cursor
2. Open Command Palette (Cmd/Ctrl + Shift + P)
3. Type "Tradovate: Get Contracts"
4. You should see the list of available contracts

### Common Issues

#### General
- Verify credentials in `.env` file
- Check file permissions (600 for `.env`)
- Ensure ports aren't conflicting (8080 for Claude, 8081 for Cursor)
- Check application logs for errors

#### Claude Desktop Specific
- Verify plugin appears in Claude's plugin manager
- Check `claude-plugin.json` syntax
- Ensure plugin directory path is correct

#### Cursor Specific
- Verify extension appears in Cursor's extension manager
- Check `cursor-extension.json` syntax
- Ensure extension directory path is correct
- Verify Cursor has required permissions

## Usage Examples

### Authentication
```go
{
    "name": "your_username",
    "password": "your_password",
    "appId": "your_app_id",
    "appVersion": "1.0",
    "cid": "your_client_id",
    "sec": "your_client_secret"
}
```

### Placing an Order
```go
{
    "accountId": 12345,
    "contractId": 67890,
    "orderType": "Limit",
    "price": 100.50,
    "quantity": 1,
    "timeInForce": "Day"
}
```

### Getting Historical Data
```go
{
    "contractId": 67890,
    "startTime": "2024-03-01T00:00:00Z",
    "endTime": "2024-03-07T00:00:00Z",
    "interval": "1h"
}
```

## Error Handling

All tools return errors in a standardized format:
```go
{
    "error": "Error message description"
}
```

## Design Principles

1. **Clean Architecture**: The project follows clean architecture principles, separating concerns into distinct layers.
2. **Interface-Driven Development**: Core functionality is defined through interfaces, enabling easy testing and flexibility.
3. **Dependency Injection**: Components are loosely coupled through dependency injection.
4. **Error Handling**: Comprehensive error handling with meaningful error messages.
5. **Type Safety**: Strong typing throughout the codebase.
6. **Documentation**: Thorough documentation of code and API endpoints.

## Smithery.ai Integration

This plugin is built to work seamlessly with Smithery.ai, following the Model Context Protocol (MCP) specification. Key integration features include:

1. **Standard MCP Interface**: Implements the MCP protocol for tool definitions and interactions
2. **Hosted Deployment**: Can be deployed and managed through Smithery.ai's platform
3. **Configuration Management**: Credentials and settings can be managed through Smithery.ai's dashboard
4. **Version Control**: Updates and versioning handled through Smithery.ai's infrastructure

## Author

Jake Peterson ([@0jxmp](https://github.com/0jxmp))

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License - see LICENSE file for details 