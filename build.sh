#!/bin/sh

if [ "$1" = "run" ]; then
    go run  -ldflags "-s -w -H=windowsgui -extldflags=-static" -race -gcflags=all=-d=checkptr=0 main.go
else
    go build  -ldflags "-s -w -H=windowsgui -extldflags=-static" -race -gcflags=all=-d=checkptr=0 main.go
fi

echo Done