package main

import (
	"flag"

	"github.com/TheLazyLemur/gofit"
	"github.com/TheLazyLemur/gofit/src/internal/server"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()

	deps := &dependencies{}
	defer deps.Close()

	dbc := deps.DBC()

	if _, err := dbc.Exec(string(gofit.Schema)); err != nil {
		panic(err)
	}

	s := server.NewServer(":8080", deps)

	server.MountRoutes(s)

	if err := server.Start(s); err != nil {
		panic(err)
	}
}
