#!/bin/bash

cd bank-test-client

go mod tidy
go build -o bank-client ./cmd/main.go

./bank-client
