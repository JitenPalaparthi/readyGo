#!/bin/sh

rm -rf grpcmongoproduct

go run main.go gen -f misc/test_configs/config_grpc_mongo_product.json -t grpc_mongo


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcmongoproduct/protos/address.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcmongoproduct/protos/product.proto

cd grpcmongoproduct

go mod init grpcmongoproduct

go run main.go



