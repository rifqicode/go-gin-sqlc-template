#!/bin/sh
# install go-migrate first
# https://github.com/golang-migrate/migrate
echo "Please enter migration name: "
read migratename

migrate create -ext sql -dir migrations -seq $migratename

