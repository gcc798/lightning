package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUpFromEmptyDatabase(t *testing.T) {
	dsn := os.Getenv("LIGHTNING_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LIGHTNING_TEST_POSTGRES_DSN is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := Up(db); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"s_config", "s_dict_data", "s_login_log", "s_oper_log"} {
		var exists bool
		if err := db.QueryRow(`SELECT to_regclass('public.' || $1) IS NOT NULL`, table).Scan(&exists); err != nil || !exists {
			t.Fatalf("table %s: exists=%v err=%v", table, exists, err)
		}
	}
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM s_config`).Scan(&count); err != nil || count != 5 {
		t.Fatalf("runtime configs=%d err=%v", count, err)
	}
}
