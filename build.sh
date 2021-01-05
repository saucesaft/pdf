#!/bin/sh

echo Starting build...

CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOARCH=amd64 GOOS=windows go build  -ldflags='-extldflags "-static"' -tags=static main.go

echo Finished
