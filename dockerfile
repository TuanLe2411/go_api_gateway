FROM golang:1.24.0-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api_gateway cmd/main/main.go

# Create a minimal production image
FROM alpine:3.14 AS final

# Add necessary runtime packages and security configurations
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/api_gateway .

# Use non-root user for better security
USER appuser

# Expose the service port
EXPOSE 8000

ENV ENV=production

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8000/health || exit 1

# Command to run the application
CMD ["./api_gateway"]