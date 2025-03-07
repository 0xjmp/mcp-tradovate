# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-tradovate ./cmd/mcp-tradovate

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and config from builder
COPY --from=builder /app/mcp-tradovate .
COPY --from=builder /app/smithery.yaml .

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Run the binary (no need to expose port as we're using STDIO)
CMD ["./mcp-tradovate"] 