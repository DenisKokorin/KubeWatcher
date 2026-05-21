#!/bin/sh
# Run SQL migrations against POSTGRES_URL env var or default
goose -dir .\migrations\ postgres "postgres://test:testpass@localhost:5432/kubewatcher_test" up
echo "Running migrations against $DB_URL"
psql "$DB_URL" -f migrations/001_init_schema.sql
psql "$DB_URL" -f migrations/002_add_role_column.sql
echo "Migrations applied."
