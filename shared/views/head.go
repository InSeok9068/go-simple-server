package shared

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HeadsWithBulma(title string) []Node {
	return []Node{
		Link(Rel("icon"), Href("/shared/static/favicon.ico")),
		Link(Rel("stylesheet"), Href("/shared/static/style.css")),
		Link(Rel("stylesheet"), Href("https://unpkg.com/open-props")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css")),
		Script(Src("https://unpkg.com/htmx.org@2.0.4")),
		Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
		Title(title),
	}
}

func HeadsWithPicoAndTailwind(title string) []Node {
	return []Node{
		Link(Rel("icon"), Href("/shared/static/favicon.ico")),
		Link(Rel("stylesheet"), Href("/shared/static/style.css")),
		Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
		Link(Rel("stylesheet"), Href("https://unpkg.com/open-props")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css")),
		// Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css")),
		// Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@yohns/picocss@2.2.10/css/pico.min.css")),
		Script(Src("https://unpkg.com/htmx.org@2.0.4")),
		Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
		Title(title),
	}
}
