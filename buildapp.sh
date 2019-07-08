#!/bin/bash

case $1 in
    'windows')
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./client.go
        zip ./client-window.zip ./client.exe ./config.json
        rm -rf ./main.exe
    ;;
    'linux')
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./client.go
        zip ./client-linux.zip ./client ./config.json
        rm -rf ./main
    ;;
    'mac')
        go build ./main.go
        zip ./client-mac.zip ./client ./config.json
        rm -rf ./main
    ;;
esac


