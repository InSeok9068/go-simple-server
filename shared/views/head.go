package shared

import (
	"simple-server/internal/config"
	"simple-server/pkg/util/gomutil"
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func HeadsDefault(title string) []Node {
	return []Node{
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		Link(Rel("icon"), Href("/shared/static/favicon.ico")),
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/open-props@1.7.16/open-props.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/open-props.min.css';")),
		Script(Src("https://cdn.jsdelivr.net/npm/htmx.org@2.0.6/dist/htmx.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/htmx.min.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/htmx-ext-alpine-morph@2.0.0/alpine-morph.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/alpine-morph.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/@alpinejs/persist@3.x.x/dist/cdn.min.js"), Defer(),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/persist.cdn.min.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js"), Defer(),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/morph.cdn.min.js';")),
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
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bulma@1.0.4/css/bulma.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/bulma.min.css';")),
	}

	return gomutil.MergeHeads(
		bulma,
		HeadsDefault(title),
		HeadsCustom(),
	)
}

func HeadsWithTailwind(title string) []Node {
	tailwind := []Node{
		Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
	}

	return gomutil.MergeHeads(
		tailwind,
		HeadsDefault(title),
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
		picoAndTailwind,
		HeadsDefault(title),
		HeadsCustom(),
	)
}

func HeadsWithBeer(title string) []Node {
	beer := []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/beercss@3.11.33/dist/cdn/beer.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/beer.min.css';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/beercss@3.11.33/dist/cdn/beer.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/beer.min.js';")),
		Script(Type("module"), Src("https://cdn.jsdelivr.net/npm/material-dynamic-colors@1.1.2/dist/cdn/material-dynamic-colors.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/material-dynamic-colors.min.js';")),
	}

	return gomutil.MergeHeads(
		beer,
		HeadsDefault(title),
		HeadsCustom(),
	)
}

func HeadWithFirebaseLogin() []Node {
	return []Node{
		Link(Rel("stylesheet"), Href("https://www.gstatic.com/firebasejs/ui/6.1.0/firebase-ui-auth.css")),
		Script(Src("https://www.gstatic.com/firebasejs/10.0.0/firebase-app-compat.js")),
		Script(Src("https://www.gstatic.com/firebasejs/10.0.0/firebase-auth-compat.js")),
		Script(Src("https://www.gstatic.com/firebasejs/ui/6.1.0/firebase-ui-auth__ko.js")),
		Script(Src("/shared/static/firebase_login.js"), Defer()),
	}
}

func HeadWithFirebaseAuth() []Node {
	return []Node{
		Script(Type("module"), Src("/shared/static/firebase_auth.js"), Defer()),
	}
}

func HeadGoogleFonts(font ...string) []Node {
	return []Node{
		Link(Rel("stylesheet"),
			Href("https://fonts.googleapis.com/css2?family="+strings.Join(font, "&family=")+"&display=swap&subset=korean")),
	}
}

func HeadFlatpickr() []Node {
	return []Node{
		Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/flatpickr.min.css"),
			Attr("onerror", "this.onerror=null;this.href='shared/static/lib/flatpickr.min.css';")),
		Script(Src("https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/flatpickr.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/flatpickr.min.js';")),
		Script(Src("https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/l10n/ko.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/ko.js';")),
	}
}

func HeadMarked() []Node {
	return []Node{
		Script(Src("https://cdn.jsdelivr.net/npm/marked/lib/marked.umd.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/marked.umd.js';"),
		),
	}
}

func HeadHammer() []Node {
	return []Node{
		Script(Src("https://cdn.jsdelivr.net/npm/hammerjs@2.0.8/hammer.min.js"),
			Attr("onerror", "this.onerror=null;this.src='shared/static/lib/hammer.min.js';")),
	}
}
