package main

import (
	"flag"

	"github.com/TheLazyLemur/gofit"
	"github.com/TheLazyLemur/gofit/src/internal/server"
	"github.com/joho/godotenv"

	_ "embed"
)

func loadEnv() {
	godotenv.Load()
}

func main() {
	flag.Parse()
	loadEnv()

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
