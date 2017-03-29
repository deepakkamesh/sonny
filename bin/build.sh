#!/bin/sh
BUILDTIME="`date '+%Y-%m-%d_%I:%M:%S%p'`"
GITHASH="`git rev-parse --short=7 HEAD`"
VER="-X main.buildtime=$BUILDTIME -X main.githash=$GITHASH"

# Compile protobufs
protoc -I ../sonny -I /usr/include/google/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "$VER" ../main.go


