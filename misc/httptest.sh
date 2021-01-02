#!/bin/sh

rm -rf mongosample

go run main.go apply -f misc/test_configs/config_mongo.json -t http_mongo

cd mongosample

go mod init mongosample

go run main.go



