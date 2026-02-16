#!/bin/bash

# Setup script for github.com/bbapp-org/auth-service
set -e

echo "Setting up Auth Service..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please update the .env file with your configuration before running the service."
fi

# Download dependencies
echo "Downloading Go dependencies..."
go mod tidy

# Run database migrations
echo "Running database migrations..."
go run scripts/migrate.go

# Seed initial data
echo "Seeding initial data..."
go run scripts/seed.go

# Generate swagger documentation
echo "Generating swagger documentation..."
swag init -g main.go -o docs/ || echo "Swagger generation failed. Install swag: go install github.com/swaggo/swag/cmd/swag@latest"

echo "Setup complete!"
echo ""
echo "To start the service:"
echo "  go run main.go"
echo ""
echo "Or use the Makefile:"
echo "  make run"
echo ""
echo "Default super admin credentials:"
echo "  Email: admin@example.com"
echo "  Password: admin123"
