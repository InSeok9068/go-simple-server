package resources

import (
	"embed"
)

//go:embed .env
//go:embed shared/static/*
//go:embed projects/sample/static/*
//go:embed projects/homepage/static/*
var EmbeddedFiles embed.FS
