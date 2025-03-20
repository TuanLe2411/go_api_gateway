FROM golang:1.24.0-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
    
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
RUN go build -ldflags="-s -w" -o api_gateway cmd/main/main.go

# Create a minimal production image
FROM alpine:3.14 AS final

# Add necessary runtime packages and security configurations
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/log && chown appuser:appgroup /app/log && chmod 777 /app/log

# Set the timezone to Asia/Ho_Chi_Minh
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/api_gateway .

COPY .env.* ./

COPY config/ ./config/

# Use non-root user for better security
USER appuser

# Expose the service port
EXPOSE 8000

ENV ENV=production

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8000/health || exit 1

# Command to run the application
ENTRYPOINT ["./api_gateway"]