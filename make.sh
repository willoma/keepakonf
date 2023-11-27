#!/bin/sh

go generate ./...
go build -o keepakonf cmd/service/main.go
upx -9 keepakonf
