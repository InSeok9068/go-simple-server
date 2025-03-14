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
		Link(Rel("stylesheet"), Href("https://unpkg.com/open-props"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/open-props.min.css';")),
		Script(Src("https://unpkg.com/htmx.org@2.0.4"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/htmx.min.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer(),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/cdn.min.js';")),
		Title(title),
	}
}

func HeadsWithBulma(title string) []Node {
	bulma := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/bulma.min.css';")),
	}

	return append(
		HeadsDefault(title),
		bulma...,
	)
}

func HeadsWithPicoAndTailwind(title string) []Node {
	picoAndTailwind := []Node{
		Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/pico.classless.min.css';")),
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
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/beercss@3.9.7/dist/cdn/beer.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/beer.min.css';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/beercss@3.9.7/dist/cdn/beer.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/beer.min.js';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/material-dynamic-colors@1.1.2/dist/cdn/material-dynamic-colors.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/material-dynamic-colors.min.js';")),
	}

	return append(
		HeadsDefault(title),
		beer...,
	)
}
