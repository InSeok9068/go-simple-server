package resources

import (
	"embed"
)

//go:embed .env
//go:embed model.conf
//go:embed shared/static/*
//go:embed projects/deario/migrations/*.sql
//go:embed projects/portfolio/migrations/*.sql
//go:embed projects/homepage/static/*
//go:embed projects/ai-study/static/*
//go:embed projects/deario/static/*
//go:embed projects/portfolio/static/*
var EmbeddedFiles embed.FS
