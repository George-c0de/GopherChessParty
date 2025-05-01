#!/usr/bin/env bash
set -e

# Ждём доступности Postgres
until pg_isready -h postgres -p 5432 -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

echo "Running migrations..."
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:postgres@postgres:5432/gopher_chess?search_path=public&sslmode=disable" \
  --revisions-schema public

echo "Starting application..."
./main --config config/dev.yaml