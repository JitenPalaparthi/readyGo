#!/bin/sh

rm -rf grpcsample

go run main.go apply -f box/configs/config_grpc_mongo.json -t grpc_mongo

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcsample/protos/address.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcsample/protos/person.proto

cd grpcsample

go mod init grpcsample

go run main.go



