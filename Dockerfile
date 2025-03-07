# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-tradovate ./cmd/mcp-tradovate

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/mcp-tradovate .

# Create non-root user
RUN adduser -D appuser
USER appuser

# Command to run the executable
ENTRYPOINT ["./mcp-tradovate"] 