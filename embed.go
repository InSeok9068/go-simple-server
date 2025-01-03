package resources

import (
	"embed"
)

//go:embed .env
//go:embed shared/static/*
//go:embed projects/sample/static/*
var EmbeddedFiles embed.FS
