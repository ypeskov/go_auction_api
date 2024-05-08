#!/usr/bin/env bash

source ../.env

user=${DB_USER:-postgres}
password=${DB_PASSWORD:-auction}
db_host=${DB_HOST:-localhost}
db_port=${DB_PORT:-5432}
db_name=${DB_NAME:-auction}

# Проверка наличия необходимых параметров
if [ -z "$user" ]; then
  echo "No database user specified."
  exit 1
fi

if [ -z "$password" ]; then
  echo "No database password specified."
  exit 1
fi

if [ -z "$db_host" ]; then
  echo "No database host specified."
  exit 1
fi

if [ -z "$db_port" ]; then
  echo "No database port specified."
  exit 1
fi

if [ -z "$db_name" ]; then
  echo "No database name specified."
  exit 1
fi

action=$1

if [ -z "$action" ]; then
  echo "No action specified: up or down."
  exit 1
fi

if [ "$action" != "up" ] && [ "$action" != "down" ]; then
  echo "Invalid action: $action."
  echo "Possible values: up, down."
  exit 1
fi

db_url="postgres://$user:$password@$db_host:$db_port/$db_name?sslmode=disable"

migrate -database "$db_url" -path migrations $action
