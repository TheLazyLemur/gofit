package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/TheLazyLemur/gofit/src/internal/db"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type dependencies struct {
	dbc     *sql.DB
	querier db.Querier
}

func (d *dependencies) DBC() *sql.DB {
	newConn := func() *sql.DB {
		dbURL := os.Getenv("GOFIT_DB_URL")
		if dbURL == "" {
			slog.Warn("GOFIT_DB_URL is not set, using in-memory database")
			dbURL = "file:memdb1?mode=memory&cache=shared"
		}

		dbc, err := sql.Open("libsql", dbURL)
		if err != nil {
			panic(err)
		}

		return dbc
	}

	if d.dbc == nil {
		d.dbc = newConn()
	}

	if err := d.dbc.Ping(); err != nil {
		d.dbc = newConn()
	}

	return d.dbc
}

func (d *dependencies) Querier() db.Querier {
	if d.querier == nil {
		d.querier = db.New()
	}

	return d.querier
}

func (d *dependencies) Close() {
	d.dbc.Close()
	d.querier = nil
	d.dbc = nil
}
