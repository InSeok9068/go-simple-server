package layout

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AppHeader() Node {
	return Header(Class("fixed primary"),
		Nav(
			A(Class("max"), Href("/"),
				H3(Text("Deario")),
			),
			Button(Class("transparent circle"),
				I(Text("search")),
			),
			A(Class("a‑login"), ID("login"), Href("/login"), x.Data(""),
				x.Show("!$store.auth.isAuthed"),
				Text("로그인"),
			),
			A(Class("a‑login"), ID("logout"), Href("#"), x.Data(""),
				x.Show("$store.auth.isAuthed"),
				Attr("onclick", "logoutUser()"),
				Text("로그아웃"),
			),
		),
	)
}
