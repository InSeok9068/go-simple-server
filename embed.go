package resources

import (
	"embed"
)

//go:embed model.conf
//go:embed shared/static/*
//go:embed projects/*/migrations/*.sql
//go:embed projects/*/static/*
var EmbeddedFiles embed.FS
