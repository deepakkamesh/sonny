#!/bin/sh
# $1 = arm or linux
# $2 binary to build: main or cli.
# $3 host to push binary to. eg. 10.0.0.20

BUILDTIME="`date '+%Y-%m-%d_%I:%M:%S%p'`"
GITHASH="`git rev-parse --short=7 HEAD`"
VER="-X main.buildtime=$BUILDTIME -X main.githash=$GITHASH"

# Fix binary paths
if [ "$2" == "main" ]; then
	BINARY="../main.go"
elif [ "$2" == "cli" ]; then
	BINARY="../controllercli/cli.go"
fi


# Compile protobuf
if [ "$(uname)" == "Darwin" ]; then
	protoc -I ../sonny -I /usr/local/include/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
	protoc -I ../sonny -I /usr/include/google/protobuf/  ../sonny/sonny.proto --go_out=plugins=grpc:../sonny
fi

# Compile binary.
if [ $1 == "arm" ]; then
	echo "Compiling for ARM"
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags "$VER" $BINARY
else
	echo "Compiling on local machine"
	go build -ldflags "$VER" $BINARY
fi

# Push binary to remote if previous step completed.
if ! [ -z "$3" ]; then
 	if [ $4 == "all" ]; then
  	echo "Pushing binary to machine $3"
		scp $2 $3:~/
  	echo "Pushing resources to machine $3"
		scp -r ../resources $3:~/
	elif [ $4 == "res" ]; then
  	echo "Pushing resources to machine $3"
		scp -r ../resources $3:~/
	elif [ $4 == "bin" ]; then
  	echo "Pushing binary to machine $3"
		scp $2 $3:~/
	fi
fi

