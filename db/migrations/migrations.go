package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var sqlFiles embed.FS

// Run applies all pending migrations.
func Run(db *sql.DB) error {
	goose.SetBaseFS(sqlFiles)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
