# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy service dependencies for proto files
COPY service-1-user/proto/ ./service-1-user/proto/
COPY service-1-user/go.mod service-1-user/go.sum ./service-1-user/
COPY service-2-article/proto/ ./service-2-article/proto/
COPY service-2-article/go.mod service-2-article/go.sum ./service-2-article/

# Copy service-3-gateway files
COPY service-3-gateway/go.mod service-3-gateway/go.sum ./service-3-gateway/

# Download dependencies
WORKDIR /build/service-3-gateway
RUN go mod download

# Copy source code
COPY service-3-gateway/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/server/

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/service-3-gateway/gateway .

# Expose HTTP port
EXPOSE 8080

# Run the application
CMD ["./gateway"]
