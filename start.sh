#!/usr/bin/env bash
set -e

# Ждём доступности Postgres
until nc -z postgres 5432; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

echo "Starting application..."
./main --config config/dev.yaml