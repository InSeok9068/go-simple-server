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
					A(ID("login"), Href("/login"), x.Data(""),
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
							aiFeedbackButton("1", "칭찬받기"),
							aiFeedbackButton("2", "위로받기"),
							aiFeedbackButton("3", "충고받기"),
							aiFeedbackButton("4", "그림일기"),
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
				A(h.Get("/diary/random"),
					lucide.Dices(),
				),
				A(h.Get("/diary/list"),
					h.Target("#diary-list-content"),
					h.On("htmx:after-on-load", "showModal('#diary-list-dialog')"),
					I(Text("list_alt")),
				),
				A(Attr("onclick", "location.reload()"),
					I(Text("refresh")),
				),
			),
			/* Footer */

			/* Dialog */
			Dialog(ID("ai-feedback-dialog"), Class("max"),
				H5(Text("일기 요정")),
				Div(ID("ai-feedback-content"),
					Text("안녕"),
				),
				Nav(Class("right-align"),
					Button(Attr("onclick", "closeModal('#ai-feedback-dialog')"),
						Text("확인"),
					),
				),
			),

			/* Dialog */
			Dialog(ID("diary-list-dialog"), Class("max"),
				H5(Text("작성 일지")),
				Ul(ID("diary-list-content"), Class("list border no-space")),
				//Nav(Class("center-align"), x.Data("{ page : 2 }"),
				//	Button(I(Text("arrow_drop_down")), h.Get("/diary/list?page=1"), h.Swap("#diary-list-content")),
				//),
				Nav(Class("right-align"),
					Button(Attr("onclick", "closeModal('#diary-list-dialog')"),
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

func aiFeedbackButton(strType string, title string) Node {
	return Li(
		h.Post(fmt.Sprintf("/ai-feedback?type=%s", strType)),
		h.Include("[name='content']"), h.Target("#ai-feedback-content"),
		h.Indicator("#feedback-loading"),
		h.On("htmx:after-on-load", "showModal('#ai-feedback-dialog')"),
		Text(title),
	)
}
