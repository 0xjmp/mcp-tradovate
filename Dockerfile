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

# Copy the binary from builder
COPY --from=builder /app/mcp-tradovate .
COPY --from=builder /app/smithery.json .

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Expose the port specified in smithery.json
EXPOSE 8080

# Run the binary
CMD ["./mcp-tradovate"] 