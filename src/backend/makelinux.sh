#!/usr/bin/bash.exe
export GOOS=linux
export GOARCH=amd64

go build -v -o ../../RollOut/Server/server server.go

