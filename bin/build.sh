#!/bin/sh
# Compile protobufs
protoc -I ../sonny -I /usr/include/google/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
go build ../main.go
