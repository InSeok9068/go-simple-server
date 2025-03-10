package shared

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HeadsDefault(title string) []Node {
	return []Node{
		Link(Rel("icon"), Href("/shared/static/favicon.ico")),
		Link(Rel("stylesheet"), Href("/shared/static/style.css")),
		Link(Rel("stylesheet"), Href("/static/style.css")),
		Link(Rel("stylesheet"), Href("https://unpkg.com/open-props")),
		Script(Src("https://unpkg.com/htmx.org@2.0.4")),
		Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
		Title(title),
	}
}

func HeadsWithBulma(title string) []Node {
	bulma := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css")),
	}

	return append(
		HeadsDefault(title),
		bulma...,
	)
}

func HeadsWithPicoAndTailwind(title string) []Node {
	picoAndTailwind := []Node{
		Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css")),
		// Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css")),
		// Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@yohns/picocss@2.2.10/css/pico.min.css")),
	}

	return append(
		HeadsDefault(title),
		picoAndTailwind...,
	)
}

func HeadsWithBeer(title string) []Node {
	beer := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/beercss@3.9.7/dist/cdn/beer.min.css")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/beercss@3.9.7/dist/cdn/beer.min.js")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/material-dynamic-colors@1.1.2/dist/cdn/material-dynamic-colors.min.js")),
	}

	return append(
		HeadsDefault(title),
		beer...,
	)
}
