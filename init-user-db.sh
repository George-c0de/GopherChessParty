#!/usr/bin/env bash
set -e

# Сначала изменяем пароль суперпользователя postgres
psql --username "postgres" <<-EOSQL
  ALTER USER postgres WITH PASSWORD '${DB_PASSWORD}';
EOSQL

# Потом создаём роль и базу данных
psql --username "postgres" <<-EOSQL
  CREATE ROLE IF NOT EXISTS ${DB_USER} LOGIN PASSWORD '${DB_PASSWORD}';
  CREATE DATABASE IF NOT EXISTS ${DB_DATABASE} OWNER ${DB_USER};
EOSQL