package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var migrationFS embed.FS

// Up 执行所有尚未应用的 PostgreSQL 迁移。
func Up(db *sql.DB) error {
	goose.SetBaseFS(migrationFS)
	goose.SetTableName("goose_iam_version")
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db, "sql"); err != nil {
		return fmt.Errorf("run IAM migrations: %w", err)
	}
	return nil
}
