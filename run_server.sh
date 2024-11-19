#!/bin/bash

cd bank-demo-app

go mod tidy
go build -o bank-server ./cmd/main.go

./bank-server
