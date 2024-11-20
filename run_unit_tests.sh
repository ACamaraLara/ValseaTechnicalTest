#!/bin/bash
cd bank-demo-app

echo "Running unit tests for bank-demo-app..."

go test -v ./...

if [ $? -eq 0 ]; then
    echo "All tests passed!"
else
    echo "Some tests failed!"
fi
