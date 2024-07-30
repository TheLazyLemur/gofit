package main

import (
	"database/sql"

	"github.com/TheLazyLemur/gofit/src/internal/db"
)

const version = "0.0.1"

type dependencies struct {
	dbc     *sql.DB
	querier db.Querier
}

func (d *dependencies) DBC() *sql.DB {
	return d.dbc
}

func (d *dependencies) Querier() db.Querier {
	return d.querier
}

func (d *dependencies) VersionChecker() string {
	return version
}
