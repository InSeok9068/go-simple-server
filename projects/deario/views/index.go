package views

import (
	// x "github.com/glsubri/gomponents-alpine"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: append(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth()...,
		),
		Body: []Node{
			/* Header */
			Header(Class("fixed responsive yellow4"),
				Nav(
					A(Href("/"), Class("max"),
						H3(Text(title)),
					),
					A(Href("/login"),
						I(Text("login")),
					),
				),
			),
			/* Header */

			/* Body */
			Main(Class("responsive"),
				Nav(
					Div(Class("max"),
						Text("3월 1일"),
					),
					I(Text("save")),
				),
				Hr(Class("medium")),
				Div(
					Class("field textarea border"),
					Textarea(Style("height : 200px")),
				),
			),
			/* Body */

			/* Footer */
			Nav(Class("bottom s"),
				A(
					I(Text("calendar_month")),
					Input(Type("date")),
				),
				A(Href("/"),
					I(Text("home")),
				),
				A(
					I(Text("settings")),
				),
			),
			/* Footer */
		},
	})
}
