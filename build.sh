#!/bin/sh

export CGO_ENABLED=0
GOOS=linux go build .
mv pask pask-linux-amd64
GOOS=darwin go build .
mv pask pask-darwin-amd64
GOOS=windows go build .
mv pask.exe pask-windows-x64.exe
