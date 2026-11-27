#!/usr/bin/env sh
set -e

mkdir -p dist
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/siminventaris-api ./cmd/api
echo "Binary siap di: backend/dist/siminventaris-api (upload ke aaPanel)"
