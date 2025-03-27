package views

import (
	h "maragu.dev/gomponents-htmx"
	"simple-server/projects/deario/db"
	"time"

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
						Text(time.Now().Format("2006년 1월 2일")),
					),
					I(Text("save")),
				),
				Hr(Class("medium")),

				Div(h.Get("/diary"),
					h.Target("#diary"),
					h.Trigger("load delay:1s"),
					h.Swap("outerHTML"),
					NewDiary(),
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

func DiaryID(id string) Node {
	return Input(Type("hidden"), ID("id"), Name("id"), Value(id))
}

func GetDiary(diary db.Diary) Node {
	return Form(ID("diary"),
		h.Post("/save"),
		h.Trigger("every 5s"),
		h.Target("#id"),
		h.Swap("outerHTML"),
		Div(ID("content"),
			DiaryID(diary.ID),
			Class("field textarea border"),
			Textarea(Name("content"), Style("height : 200px"),
				Text(diary.Content),
			),
		),
	)
}

func NewDiary() Node {
	return Form(ID("diary"),
		h.Post("/save"),
		h.Trigger("every 5s"),
		h.Target("#id"),
		h.Swap("outerHTML"),
		Div(ID("content"),
			DiaryID(""),
			Class("field textarea border"),
			Textarea(Name("content"), Style("height : 200px")),
		),
	)
}
