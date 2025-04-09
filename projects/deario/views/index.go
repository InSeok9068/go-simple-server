package views

import (
	"fmt"
	lucide "github.com/eduardolat/gomponents-lucide"
	x "github.com/glsubri/gomponents-alpine"
	h "maragu.dev/gomponents-htmx"
	"simple-server/pkg/util"
	shared "simple-server/shared/views"
	"time"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string, date string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: util.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("/static/deario.js")),
			},
		),
		Body: []Node{
			/* Header */
			Header(Class("fixed responsive yellow4"),
				Nav(
					A(Href("/"), Class("max"),
						H3(Text(title)),
					),
					A(ID("login"), Href("/login"), x.Data(""), /*, x.Show("!$store.auth.isAuthed")*/
						Text("Login"),
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
					A(Href(fmt.Sprintf("/?date=%s", time.Now().Format("20060102"))),
						Button(Class("small-round small"), Text("Today")),
					),
					A(Href(fmt.Sprintf("/?date=%s", util.MustAddDaysToDate(date, 1))),
						I(Text("arrow_right")),
					),
					Div(Class("max"),
						Text(DateView(date)),
					),
					P(Class("max bold"), x.Data(""), x.Show("$store.auth.isAuthed"), x.Text("$store.auth?.user?.displayName")),
				),
				Hr(Class("medium")),

				// 일기장 조회
				Form(h.Get("/diary"), h.Target("#diary"), h.Trigger("load"), h.Swap("outerHTML"),
					Input(Type("hidden"), Name("date"), Value(date)),
				),

				DiaryContentForm(date, ""),

				Nav(
					Button(Class("chip circle"),
						I(Text("psychology_alt")),
						Div(Class("tooltip right"),
							Ul(Class("list no-space"),
								Li(Text("- 짧게라도 하루를 되돌아보기")),
								Li(Text("- 감정을 솔직하게 작성")),
							),
						),
					),
					Div(Class("max")),
					Img(ID("feedback-loading"), Class("htmx-indicator"), Src("/shared/static/spinner.svg")),
					Button(Attr("data-ui", "#ai-feedback"),
						Span(
							Text("일기 요정"),
						),
						Menu(Class("top no-wrap"), ID("ai-feedback"), Attr("data-ui", "#ai-feedback"),
							Li(h.Post("/ai-feedback?type=1"), h.Include("[name='content']"), h.Target("#ai-feedback-content"),
								h.Indicator("#feedback-loading"),
								h.On("htmx:after-on-load", "document.querySelector('#ai-feedback-modal').showModal()"),
								Text("칭찬받기"),
							),
							Li(h.Post("/ai-feedback?type=2"), h.Include("[name='content']"), h.Target("#ai-feedback-content"),
								h.Indicator("#feedback-loading"),
								h.On("htmx:after-on-load", "document.querySelector('#ai-feedback-modal').showModal()"),
								Text("위로받기"),
							),
							Li(h.Post("/ai-feedback?type=3"), h.Include("[name='content']"), h.Target("#ai-feedback-content"),
								h.Indicator("#feedback-loading"),
								h.On("htmx:after-on-load", "document.querySelector('#ai-feedback-modal').showModal()"),
								Text("충고받기"),
							),
							Li(h.Post("/ai-feedback?type=4"), h.Include("[name='content']"), h.Target("#ai-feedback-content"),
								h.Indicator("#feedback-loading"),
								h.On("htmx:after-on-load", "document.querySelector('#ai-feedback-modal').showModal()"),
								Text("그림일기"),
							),
						),
					),
				),
			),
			/* Body */

			/* Footer */
			Nav(Class("bottom"),
				A(
					I(Text("calendar_month")),
					Input(Type("date"), Name("date"), Attr("onchange", "location.href='/?date=' + this.value")),
				),
				A(Attr("onclick", "location.reload()"),
					I(Text("refresh")),
				),
				A(h.Get("/diary/random"),
					lucide.Dices(),
				),
			),
			/* Footer */

			/* Dialog */
			Dialog(ID("ai-feedback-modal"), Class("max"),
				H5(Text("일기 요정")),
				Div(ID("ai-feedback-content"),
					Text("안녕"),
				),
				Nav(Class("right-align"),
					Button(Attr("onclick", "document.querySelector('#ai-feedback-modal').close()"),
						Text("확인"),
					),
				),
			),
		},
	})
}

func DateView(date string) string {
	parsed, _ := time.Parse("20060102", date)
	dateStr := parsed.Format("1월 2일")
	return dateStr
}

func DiaryContentForm(date string, content string) Node {
	return Form(ID("diary"),
		Input(Type("hidden"), Name("date"), Value(date)),
		Div(Class("field textarea border medium-height"),
			Textarea(
				Name("content"),
				AutoFocus(),
				h.Post("/save"),
				h.Swap("none"),
				h.Trigger("input delay:0.5s"),
				h.Indicator("#indicator"),
				Text(content),
			),
			Img(ID("indicator"), Class("htmx-indicator"), Src("/shared/static/spinner.svg")),
		),
	)
}
