#!/bin/sh

MIGRATION_DIR="./migrations"
DB_DSN="host=localhost port=5432 user=test password=test dbname=test sslmode=disable"

if [ "$1" = "--dryrun" ]; then
    ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" status
elif [ "$1" = "--down" ]; then
  ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" down
else
    ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" up
fi