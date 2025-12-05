# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go.mod files
COPY go.mod go.sum ./

# Copy vendor directory (contains all dependencies including proto files)
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Build the application with vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o gateway ./cmd/server/

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gateway .

# Expose HTTP port
EXPOSE 8080

# Run the application
CMD ["./gateway"]
