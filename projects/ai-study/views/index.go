package views

import (
	x "github.com/glsubri/gomponents-alpine"
	b "github.com/willoma/bulma-gomponents"
	e "github.com/willoma/gomplements"
	. "maragu.dev/gomponents"
	h "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	shared "simple-server/shared/views"
)

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head:     shared.HeadsWithBulma(title),
		Body: []Node{
			b.Container(b.MaxTablet, b.Padding(2),
				b.Columns(
					b.Column(
						b.Box(
							b.Content(
								H3(Raw("오늘 무엇을 공부하고싶은지 <br> 고민이신가요 ?")),
								P(Text("AI와 함께 공부할 수 있도록 주제와 사이트 접속 링크를 제공해드려요")),
								Form(
									h.Post("/ai-study"),
									h.Target("#result"),
									x.Data("{ loading : false }"),
									x.On("htmx:before-request", "loading = true"),
									x.On("htmx:after-on-load", "loading = false"),
									b.Field(b.Addons,
										b.Control(b.Expanded,
											b.InputText(Name("input"), Placeholder("내용을 입력해주세요."), Required()),
										),
										b.Control(
											b.Button(b.Info,
												Text("입력"),
												Img(x.Show("loading"), Class("htmx-indicator"), Src("/static/spinner.svg")),
											),
										),
									),
								),
								b.Columns(
									b.Column(b.Half, b.OffsetOneQuarter,
										b.Button(b.Focused, b.FullWidth, b.MarginTop(3),
											Type("submit"),
											h.Post("/ai-study-random"),
											h.Target("#result"),
											x.Data("{ loading : false }"),
											x.On("htmx:before-request", "loading = true"),
											x.On("htmx:after-on-load", "loading = false"),
											Text("❔ 랜덤 주제로 시작하기"),
											Img(x.Show("loading"), Class("htmx-indicator"), Src("/static/spinner.svg")),
										),
									),
								),
							),
						),
						P(Text("⏬ 결과 확인")),
						b.Box(e.Class("content"),
							b.Columns(b.Mobile,
								b.Column(
									b.Button(ID("copy"),
										Text("글 복사")),
								),
								b.Column(
									A(Href("https://chatgpt.com/"), Target("_blank"),
										Text("ChatGPT")),
								),
								b.Column(
									A(Href("https://gemini.google.com/app?hl=ko"), Target("_blank"),
										Text("Gemini")),
								),
							),
							Div(ID("result")),
						),
					),
				),
			),
			Script(Raw(`
				document.getElementById("copy").addEventListener("click", () => {
					navigator.clipboard.writeText(
						document.getElementById("result").innerText
					).then(() => alert("복사 되었습니다."));
				});
			`)),
		},
	})
}
