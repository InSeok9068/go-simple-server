package views

import (
	"fmt"
	x "github.com/glsubri/gomponents-alpine"
	h "maragu.dev/gomponents-htmx"
	"simple-server/pkg/util"
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
					A(ID("login"), Href("/login"), x.Data(""), x.Show("!$store.auth.isAuthed"),
						I(Text("login")),
					),
					A(ID("logout"), x.Data(""), x.Show("$store.auth.isAuthed"),
						I(Text("logout")),
					),
				),
			),
			/* Header */

			/* Body */
			Main(Class("responsive"),
				Nav(
					A(Href(fmt.Sprintf("/?date=%s", util.MustAddDaysToDate(date, -1))),
						I(Text("arrow_left")),
					),
					A(Href(fmt.Sprintf("/?date=%s", util.MustAddDaysToDate(date, 1))),
						I(Text("arrow_right")),
					),
					Div(Class("max"),
						Text(DateView(date)),
					),
					P(Class("bold"), x.Data(""), x.Show("$store.auth.isAuthed"), x.Text("$store.auth?.user?.displayName")),
				),
				Hr(Class("medium")),

				// 일기장 조회
				Form(h.Get("/diary"), h.Target("#diary"), h.Trigger("firebase:authed"), h.Swap("outerHTML"),
					Input(Type("hidden"), Name("date"), Value(date)),
				),

				DiaryContentForm(date, ""),

				Nav(
					Div(Class("max")),
					Button(Attr("data-ui", "#ai-feedback"),
						Span(Text("AI 피드백")),
						Menu(Class("top no-wrap"), ID("ai-feedback"), Attr("data-ui", "#ai-feedback"),
							Li(Text("칭찬받기 ^^")),
							Li(Text("위로받기 ㅜㅜ")),
							Li(Text("충고받기 !!")),
						),
					),
				),
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

func DiaryContentForm(date string, content string) Node {
	return Form(ID("diary"),
		Input(Type("hidden"), Name("date"), Value(date)),
		Div(Class("field textarea border medium-height"),
			Textarea(
				Name("content"),
				h.Post("/save"),
				h.Swap("none"),
				h.Trigger("input delay:0.5s"),
				h.Indicator("#indicator"),
				Style("height : 400px"),
				Text(content),
			),
			Img(ID("indicator"), Class("htmx-indicator"), Src("/shared/static/spinner.svg")),
		),
	)
}
