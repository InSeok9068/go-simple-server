package views

import (
	h "maragu.dev/gomponents-htmx"
	"simple-server/projects/deario/db"
	"time"

	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string, date string) Node {
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
						Text(DateView(date)),
					),
				),
				Hr(Class("medium")),

				// 일기장 조회
				Form(h.Get("/diary"), h.Target("#diary"), h.Trigger("firebase:authed"), h.Swap("outerHTML"),
					Input(Type("hidden"), Name("date"), Value(date)),
				),

				NewDiaryContent(date),
			),
			/* Body */

			/* Footer */
			Nav(Class("bottom s"),
				A(
					I(Text("calendar_month")),
					Input(Type("date"), Name("date"), Attr("onchange", "location.href='/?date=' + this.value")),
				),
				A(Href(""),
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

func DateView(date string) string {
	parsed, _ := time.Parse("20060102", date)
	dateStr := parsed.Format("2006년 1월 2일")
	return dateStr
}

func GetDiaryContent(diary db.Diary) Node {
	return Form(ID("diary"),
		Input(Type("hidden"), Name("date"), Value(diary.Date)),
		Div(Class("field textarea border"),
			Textarea(Name("content"), h.Post("/save"), h.Swap("none"), h.Trigger("input delay:0.5s"), Style("height : 200px"),
				Text(diary.Content),
			),
		),
	)
}

func NewDiaryContent(date string) Node {
	return Form(ID("diary"),
		Input(Type("hidden"), Name("date"), Value(date)),
		Div(Class("field textarea border"),
			Textarea(Name("content"), h.Post("/save"), h.Swap("none"), h.Trigger("input delay:0.5s"), Style("height : 200px")),
		),
	)
}
