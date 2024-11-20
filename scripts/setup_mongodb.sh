#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Docker is installed
if ! command_exists docker; then
    echo "Docker is not installed. You need to install docker on your machine."
    exit 1
else
    echo "Docker is already installed."
fi

# Login into your docker account.
docker login

# Pull the MongoDB Docker image
echo "Pulling the MongoDB Docker image..."
docker pull mongo

# Create an initialization script for MongoDB
INIT_SCRIPT="mongo-init.js"
cat <<EOF > $INIT_SCRIPT
db = db.getSiblingDB('BankStore');
db.createCollection('accounts');
db.createCollection('transactions');
EOF

echo "Initialization script for MongoDB created: $INIT_SCRIPT"

# Run MongoDB container with the initialization script
echo "Starting MongoDB container..."
docker run -d \
    --name mongodb \
    -p 27017:27017 \
    -v $(pwd)/$INIT_SCRIPT:/docker-entrypoint-initdb.d/$INIT_SCRIPT:ro \
    mongo

echo "MongoDB is running with the 'BankStore' database and collections 'accounts' and 'transactions' initialized."
