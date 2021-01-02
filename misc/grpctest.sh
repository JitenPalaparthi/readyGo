#!/bin/sh

rm -rf productsample

go run main.go apply -f box/configs/config_grpc_mongo_product.json -t grpc_mongo


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative productsample/protos/address.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative productsample/protos/product.proto

cd productsample

go mod init productsample

go run main.go



