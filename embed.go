package gofit

import "embed"

//go:embed db/schema.sql
var Schema []byte

//go:embed static
var Static embed.FS
