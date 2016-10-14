#!/bin/sh

# Compile protobuf
protoc -I ../sonny -I /usr/include/google/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
# Compile main
GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build ../controllercli/cli.go

