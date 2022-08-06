#!/bin/sh

MIGRATION_DIR="./migrations"
GENERATIONS_DIR="./migrations/generations"
DB_DSN="host=localhost port=5432 user=user password=password dbname=price-tracker sslmode=disable"

if [ "$1" = "--dryrun" ]; then
    ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" status
elif [ "$1" = "--down" ]; then
  ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" down
elif [ "$1" = "--generate" ]; then
  if [ "$2" = "--down" ]; then
    ./bin/goose -v -dir ${GENERATIONS_DIR} postgres "${DB_DSN}" down
  else
    ./bin/goose -v -dir ${GENERATIONS_DIR} postgres "${DB_DSN}" up
  fi
else
    ./bin/goose -v -dir ${MIGRATION_DIR} postgres "${DB_DSN}" up
fi