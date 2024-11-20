#!/bin/bash

# Function to clear collections using the provided MongoDB shell command
clear_collections() {
  local shell_cmd=$1
  $shell_cmd "localhost:27017/BankStore" --eval "
    db.accounts.deleteMany({});
    db.transactions.deleteMany({});
    print('Collections cleared successfully.');
  " || {
    echo "Failed to clear collections using $shell_cmd."
    exit 1
  }
}

echo "Starting collection clearing process..."

# Try using mongosh first, fallback to mongo if needed
if command -v mongosh &> /dev/null; then
  clear_collections "mongosh"
elif command -v mongo &> /dev/null; then
  clear_collections "mongo"
else
  echo "Neither 'mongosh' nor 'mongo' is installed. Please install one of them to proceed."
  exit 1
fi

echo "Operation completed successfully."
