#!/usr/bin/env bash
set -e

# Потом создаём роль и базу данных
psql --username "postgres" <<-EOSQL
  CREATE ROLE IF NOT EXISTS ${DB_USER} LOGIN PASSWORD '${DB_PASSWORD}';
  CREATE DATABASE IF NOT EXISTS ${DB_DATABASE} OWNER ${DB_USER};
EOSQL