package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUpFromEmptyDatabase(t *testing.T) {
	dsn := os.Getenv("MS_K_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("MS_K_TEST_POSTGRES_DSN is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := Up(db); err != nil {
		t.Fatal(err)
	}
	var exists bool
	if err := db.QueryRow(`SELECT to_regclass('public.biz_attachment') IS NOT NULL`).Scan(&exists); err != nil || !exists {
		t.Fatalf("exists=%v err=%v", exists, err)
	}
}
