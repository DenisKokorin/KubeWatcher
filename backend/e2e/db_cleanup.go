package e2e

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// ResetDatabase truncates all public tables to provide a clean state between tests.
func ResetDatabase() error {
	dbUrl := os.Getenv("POSTGRES_URL")
	if dbUrl == "" {
		dbUrl = os.Getenv("BACKEND_POSTGRES_URL")
	}
	if dbUrl == "" {
		dbUrl = "postgres://test:testpass@localhost:5432/kubewatcher_test?sslmode=disable"
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer conn.Close(ctx)

	// Truncate all tables in public schema
	_, err = conn.Exec(ctx, `DO $$ DECLARE r RECORD; BEGIN FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE'; END LOOP; END $$;`)
	if err != nil {
		return fmt.Errorf("truncate tables: %w", err)
	}
	return nil
}
