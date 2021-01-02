
#!/bin/sh

rm -rf mysqlsample

go run main.go apply -f misc/test_configs/config_http_sql_mysql.json -t http_sql_mysql

cd mysqlsample

go mod init mysqlsample

go run main.go