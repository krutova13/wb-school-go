#!/bin/bash

set -e

if [ -f /.env ]; then
    export $(cat /.env | grep -v '^#' | xargs)
fi

DATABASE_URL="postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST_DOCKER}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}"

echo "Waiting for PostgreSQL to be ready..."
sleep 5

echo "PostgreSQL is ready - running migrations..."
echo "Database URL: $DATABASE_URL"

/migrate -path /migrations -database "$DATABASE_URL" up

echo "Migrations completed successfully"
