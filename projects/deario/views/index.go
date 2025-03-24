package views

import (
	h "maragu.dev/gomponents-htmx"
	"simple-server/projects/deario/db"
	// x "github.com/glsubri/gomponents-alpine"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string, diary *db.Diary) Node {
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
				Diary(diary),
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

func Diary(diary *db.Diary) Node {
	return Div(ID("content"),
		h.Post("/diary"),
		h.Target("#content"),
		h.Trigger("load delay:1s"),
		Iff(diary != nil, func() Node {
			return Input(Type("hidden"), ID("id"), Value(diary.ID))
		}),
		Class("field textarea border"),
		Textarea(Style("height : 200px"),
			Iff(diary != nil, func() Node {
				return Text(diary.Content)
			}),
		),
	)
}
