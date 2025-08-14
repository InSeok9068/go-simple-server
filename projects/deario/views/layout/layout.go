package layout

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	h "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func AppHeader() Node {
	return Group([]Node{
		Header(Class("fixed primary"),
			Nav(
				A(Class("max"), Href("/"),
					H3(Text("Deario")),
				),
				Button(Class("transparent circle"), Attr("data-ui", "#search-dialog"),
					I(Text("search")),
				),
				A(Href("/privacy"),
					Text("개인정보"),
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
		),
		Dialog(ID("search-dialog"),
			H5(Text("검색")),
			Form(h.Get("/diary/search"), h.Target("#search-result"), h.Swap("innerHTML"),
				Div(Class("field large prefix round fill"),
					I(Class("front"), Text("search")),
					Input(Type("search"), Name("q"), Placeholder("검색어 입력")),
				),
				Button(Type("submit"), Text("검색")),
			),
			Ul(ID("search-result"), Class("list border")),
			Nav(Class("right-align"),
				Button(Attr("data-ui", "#search-dialog"), Text("닫기")),
			),
		),
	})
}
