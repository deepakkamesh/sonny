#!/bin/sh

# Compile protobuf
protoc -I ../sonny -I /usr/include/google/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
# Compile main
go build ../controllercli/cli.go

