package main

import (
	"database/sql"
	"flag"

	"github.com/TheLazyLemur/gofit"
	"github.com/TheLazyLemur/gofit/src/internal/db"
	"github.com/TheLazyLemur/gofit/src/internal/server"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()

	dbc, err := sql.Open("sqlite3", "./file.db")
	if err != nil {
		panic(err)
	}
	defer dbc.Close()

	if _, err := dbc.Exec(string(gofit.Schema)); err != nil {
		panic(err)
	}

	deps := &dependencies{
		dbc:     dbc,
		querier: db.New(),
	}

	s := server.NewServer(":8080", deps)

	server.MountRoutes(s)

	if err := server.Start(s); err != nil {
		panic(err)
	}
}
