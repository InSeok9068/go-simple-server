package shared

import (
	"simple-server/internal/config"
	"simple-server/pkg/util/gomutil"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HeadsDefault(title string) []Node {
	return []Node{
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0, maximum-scale=1.0")),
		Link(Rel("icon"), Href("/shared/static/favicon.ico")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/open-props@1.7.13/open-props.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/open-props.min.css';")),
		Script(Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.6/dist/htmx.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/htmx.min.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"), Defer(),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/cdn.min.js';")),
		If(config.IsDevEnv(), Script(Text(`htmx.logAll();`), Defer())),
		Title(title),
	}
}

func HeadsCustom() []Node {
	custom := []Node{
		Link(Rel("stylesheet"), Href("/shared/static/style.css")),
		Link(Rel("stylesheet"), Href("/static/style.css")),
		Script(Src("/shared/static/app.js")),
	}

	return custom
}

func HeadsWithBulma(title string) []Node {
	bulma := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bulma@1.0.2/css/bulma.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/bulma.min.css';")),
	}

	return gomutil.MergeHeads(
		HeadsDefault(title),
		bulma,
		HeadsCustom(),
	)
}

func HeadsWithTailwind(title string) []Node {
	tailwind := []Node{
		Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
	}

	return gomutil.MergeHeads(
		HeadsDefault(title),
		tailwind,
		HeadsCustom(),
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

	return gomutil.MergeHeads(
		HeadsDefault(title),
		picoAndTailwind,
		HeadsCustom(),
	)
}

func HeadsWithBeer(title string) []Node {
	beer := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/beercss@3.11.32/dist/cdn/beer.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/beer.min.css';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/beercss@3.11.32/dist/cdn/beer.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/beer.min.js';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/material-dynamic-colors@1.1.2/dist/cdn/material-dynamic-colors.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/material-dynamic-colors.min.js';")),
	}

	return gomutil.MergeHeads(
		HeadsDefault(title),
		beer,
		HeadsCustom(),
	)
}

func HeadWithFirebaseLogin() []Node {
	return []Node{
		Link(Rel("stylesheet"), Href("https://www.gstatic.com/firebasejs/ui/6.1.0/firebase-ui-auth.css")),
		Script(Src("https://www.gstatic.com/firebasejs/10.0.0/firebase-app-compat.js")),
		Script(Src("https://www.gstatic.com/firebasejs/10.0.0/firebase-auth-compat.js")),
		Script(Src("https://www.gstatic.com/firebasejs/ui/6.1.0/firebase-ui-auth.js")),
		Script(Src("/shared/static/firebase_login.js"), Defer()),
	}
}

func HeadWithFirebaseAuth() []Node {
	return []Node{
		Script(Type("module"), Src("/shared/static/firebase_auth.js"), Defer()),
	}
}
