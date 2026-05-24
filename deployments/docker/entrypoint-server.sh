#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
for i in $(seq 1 30); do
  if pg_isready -h postgres -U postgres 2>/dev/null; then
    echo "PostgreSQL is ready"
    break
  fi
  echo "Waiting... ($i/30)"
  sleep 1
done

echo "Running database migrations..."
./migrate-bin -direction up || echo "Migrations already applied or failed, continuing..."

echo "Starting server..."
exec ./server
