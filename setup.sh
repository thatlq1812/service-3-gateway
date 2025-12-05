#!/bin/bash

# Gateway Service Setup Script
# This script sets up the API Gateway development environment

set -e

echo "=========================================="
echo "API Gateway Service Setup"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    echo "Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed${NC}"
    echo "Please install Docker Compose from https://docs.docker.com/compose/install/"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Warning: Go is not installed${NC}"
    echo "Go is required for local development"
    echo "Install from https://golang.org/dl/"
fi

echo ""
echo "1. Creating environment file..."
if [ ! -f .env ]; then
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${GREEN}✓ Created .env from .env.example${NC}"
    else
        echo -e "${YELLOW}Warning: .env.example not found, skipping${NC}"
    fi
else
    echo -e "${YELLOW}✓ .env already exists${NC}"
fi

echo ""
echo "2. Pulling Docker images..."
docker-compose pull

echo ""
echo "3. Building Gateway service..."
docker-compose build gateway

echo ""
echo "4. Starting services..."
docker-compose up -d

echo ""
echo "5. Waiting for services to be healthy (30 seconds)..."
sleep 30

echo ""
echo "6. Checking service status..."
docker-compose ps

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo -e "${GREEN}Gateway is running at: http://localhost:8080${NC}"
echo ""
echo "Quick commands:"
echo "  Health check:  curl http://localhost:8080/health"
echo "  List users:    curl http://localhost:8080/api/v1/users"
echo "  View logs:     docker logs gateway-app --tail 20"
echo "  Stop services: docker-compose down"
echo ""
echo "Backend services:"
echo "  User Service:     localhost:50051 (gRPC)"
echo "  Article Service:  localhost:50052 (gRPC)"
echo "  PostgreSQL User:  localhost:5434"
echo "  PostgreSQL Article: localhost:5435"
echo "  Redis:            localhost:6381"
echo ""
echo "See README.md for API documentation and testing examples"
