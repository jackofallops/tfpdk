package commands

import (
	"embed"
)

//go:embed templates/*
var Templatedir embed.FS
