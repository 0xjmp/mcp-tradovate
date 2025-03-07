# MCP Tradovate Server

![](https://badge.mcpx.dev?type=server 'MCP Server')
[![smithery badge](https://smithery.ai/badge/@0xjmp/mcp-tradovate)](https://smithery.ai/server/@0xjmp/mcp-tradovate)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xjmp/mcp-tradovate)](https://goreportcard.com/report/github.com/0xjmp/mcp-tradovate)
[![codecov](https://codecov.io/gh/0xjmp/mcp-tradovate/branch/main/graph/badge.svg)](https://codecov.io/gh/0xjmp/mcp-tradovate)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xjmp/mcp-tradovate.svg)](https://pkg.go.dev/github.com/0xjmp/mcp-tradovate)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Model Context Protocol (MCP) server for Tradovate integration in Claude Desktop. This server enables AI assistants to manage Tradovate trading accounts through natural language interactions.

## Features

- ✅ Complete Tradovate API integration
- 🔒 Secure authentication handling
- 📈 Real-time market data access
- 💼 Account management
- 📊 Risk management controls
- 🔄 Order placement and management
- 📝 Comprehensive test coverage

## Installation

### Installing via Smithery
To install the Tradovate MCP server for Claude Desktop automatically via Smithery:

```bash
npx -y @smithery/cli install @0xjmp/mcp-tradovate --client claude
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/0xjmp/mcp-tradovate.git
cd mcp-tradovate
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build ./cmd/mcp-tradovate
```

## Configuration

Create a `.env` file in the project root with your Tradovate credentials:

```env
TRADOVATE_USERNAME=your_username
TRADOVATE_PASSWORD=your_password
TRADOVATE_APP_ID=your_app_id
TRADOVATE_APP_VERSION=your_app_version
TRADOVATE_CID=your_client_id
TRADOVATE_SEC=your_client_secret
```

## Usage

1. Start the MCP server:
```bash
./mcp-tradovate
```

2. The server will listen for MCP requests on standard input and output responses on standard output.

3. Available commands:
- `authenticate`: Connect to Tradovate API
- `getAccounts`: List trading accounts
- `getPositions`: View current positions
- `placeOrder`: Submit new orders
- `setRiskLimits`: Configure risk management settings
- And more...

## Development

### Running Tests

Run all tests with coverage:
```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### Code Style

Follow Go best practices and conventions:
```bash
go fmt ./...
go vet ./...
```

## Troubleshooting

### Common Issues

1. **Authentication Failures**
   - Verify your Tradovate credentials in the `.env` file
   - Ensure your API access is enabled in Tradovate

2. **Connection Issues**
   - Check your internet connection
   - Verify Tradovate API status
   - Ensure firewall isn't blocking connections

3. **Rate Limiting**
   - Implement appropriate delays between requests
   - Monitor API usage limits

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.

## Author

Jake Peterson ([@0xjmp](https://github.com/0xjmp))

If this library helped you, consider donating (send whatever crypto you want): `0xB5BaA3D2056be942a9F61Cc015b83562DA3C15B2` 