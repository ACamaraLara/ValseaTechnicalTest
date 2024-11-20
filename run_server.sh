#!/bin/bash

# Navigate to the project directory
cd bank-demo-app

# Ensure dependencies are tidy and build the application
go mod tidy
go build -o bank-server ./cmd/main.go

# Change to false to use database instead of builtin account manager.
IN_MEMORY=false

# Execute the application based on the in-memory parameter
if [ "$IN_MEMORY" == "true" ]; then
    echo "Running in in-memory mode..."
    ./bank-server
else
    echo "Running with database configuration..."
    ./bank-server --in-memory=false --mongo-db="BankStore"
fi