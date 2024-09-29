#!/bin/sh
# install go-migrate first
# https://github.com/golang-migrate/migrate
set -e

if test -z "$POSTGRESQL_URL"; then
    echo "Please set first for \"POSTGRESQL_URL\" variable, e.g. export POSTGRESQL_URL=postgresql://localhost:5432/mydb?sslmode=disable"
    exit 1
fi

migrate -database ${POSTGRESQL_URL} -path migrations down 1

