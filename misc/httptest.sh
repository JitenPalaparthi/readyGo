#!/bin/sh

rm -rf contacts

go run main.go apply -f misc/test_configs/config_http_mongo.json -t http_mongo

cd contacts

go mod init contacts

go run main.go



