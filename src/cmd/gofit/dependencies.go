package main

import (
	"database/sql"

	"github.com/TheLazyLemur/gofit/src/internal/db"
)

type dependencies struct {
	dbc     *sql.DB
	querier db.Querier
}

func (d *dependencies) DBC() *sql.DB {
	newConn := func() *sql.DB {
		// dbc, err := sql.Open("sqlite3", "file:memdb1?mode=memory&cache=shared")
		dbc, err := sql.Open("sqlite3", "file.db")
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
