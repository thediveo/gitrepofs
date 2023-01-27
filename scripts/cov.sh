#!/bin/bash
set -e

if ! command -v go-acc; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v go-acc; then
        go install github.com/ory/go-acc@latest
    fi
fi

if ! command -v gobadge &>/dev/null; then
    export PATH="$(go env GOPATH)/bin:$PATH"
    if ! command -v gobadge &>/dev/null; then
        go install github.com/AlexBeauchemin/gobadge@latest
    fi
fi

go-acc --covermode atomic -o coverage.out ./... -- -v -p 1
go tool cover -html=coverage.out -o=coverage.html
go tool cover -func=coverage.out -o=coverage.out
gobadge -filename=coverage.out -green=80 -yellow=50
