#!/bin/bash
VERSION=v0.0.0

echo $VERSION
echo linux64
env GOOS=linux GOARCH=amd64 go build -o "bin/logcar_${VERSION}_linux_x86_64" ./
env GOOS=linux GOARCH=arm64 go build -o bin/logcar_${VERSION}_linux_arm64 ./
echo darwin
env GOOS=darwin GOARCH=arm64 go build -o bin/logcar_${VERSION}_darwin_arm64 ./
env GOOS=darwin GOARCH=amd64 go build -o bin/logcar_${VERSION}_darwin_x86_64 ./
