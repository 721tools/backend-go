package migrate

import (
	"database/sql"

	_ "github.com/721tools/backend-go/index/migration"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oh-go/goose"
)

func Migrate(dsn string, direction string, args ...string) (err error) {
	if direction == "down" {
		return rollback(dsn, args...)
	}
	return upgrade(dsn)
}

func upgrade(dsn string) (err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer db.Close()
	if err = goose.SetDialect("mysql"); err != nil {
		return
	}
	return goose.Run("up", db, ".")
}

func rollback(dsn string, args ...string) (err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer db.Close()
	if err = goose.SetDialect("mysql"); err != nil {
		return
	}
	if len(args) > 0 {
		return goose.Run("down-to", db, ".", args...)
	}
	return goose.Run("down", db, ".")
}
