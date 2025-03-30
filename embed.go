package resources

import (
	"embed"
)

//go:embed .env
//go:embed shared/static/*
//go:embed projects/homepage/static/*
//go:embed projects/ai-study/static/*
//go:embed projects/deario/static/*
var EmbeddedFiles embed.FS
