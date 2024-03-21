package database

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
)

//go:embed sql/schema.sql
var createTablesQueries string

func SetupTest() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
	}
	setup(db)
	return db
}

func Setup() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Could not open database: %s", err)
	}
	setup(db)
	return db
}

func setup(db *sql.DB) {
	ctx := context.Background()

	// create tables
	if _, err := db.ExecContext(ctx, createTablesQueries); err != nil {
		//TODO: Handle better
		log.Fatalf("Could not create tables: %s", err)
	}

}
