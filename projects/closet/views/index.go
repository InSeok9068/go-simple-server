package views

import (
	"simple-server/pkg/util/gomutil"
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Meta(Name("description"), Content("옷장 서비스")),
				Link(Rel("manifest"), Href("/manifest.json")),
			},
		),
		Body: []Node{
			shared.Snackbar(),

			/* Header */
			Header(Class("fixed primary"),
				Nav(
					A(Class("max"), Href("/"),
						H3(Text("Closet")),
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
			/* Header */

			/* Body */
			Main(ID("closet-main"),
				Div(Text("안녕하세요")),
			),
			/* Body */

			/* Footer */
			Nav(Class("bottom")),
			/* Footer */
		},
	})
}
