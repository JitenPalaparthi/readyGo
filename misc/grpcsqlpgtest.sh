#!/bin/sh

rm -rf grpcsqlpgsample

go run main.go apply -f misc/test_configs/config_grpc_sql_pg.json -t grpc_sql_pg


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcsqlpgsample/protos/address.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpcsqlpgsample/protos/product.proto

cd grpcsqlpgsample

go mod init grpcsqlpgsample

go run main.go

