package views

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	shared "simple-server/shared/views"
)

func ServiceCard(name string, desc string, url string) Node {
	return Article(
		H6(Text(name)),
		P(Text(desc)),
		Hr(Class("small")),
		Nav(
			A(Href(url), Class("underline"),
				I(
					Text("link"),
				),
				Text("서비스 이동"),
			),
		),
	)
}

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head:     shared.HeadsWithBeer(title),
		Body: []Node{
			Div(x.Data(`{ open : false }`),
				Nav(Class("bottom s"),
					A(Href("/"),
						I(Text("home")),
					),
					A(
						I(Text("search")),
					),
					A(
						I(Text("share")),
					),
					A(x.On("click", "open = !open"), Attr("data-ui", "#dialog-right"),
						I(Text("more_vert")),
					),
				),

				Dialog(ID("dialog-right"), Class("right no-padding"),
					x.Class("{ 'active': open }"),
					Nav(Class("drawer"),
						A(Href("/"), I(Text("home"))),
					),
				),

				Header(Class("fixed responsive yellow4"),
					A(Href("/"),
						H3(Text(title)),
					),
				),

				Div(Class("space")),

				Main(Class("responsive"),
					ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져줘요", "https://ai-study.toy-project.n-e.kr"),
					ServiceCard("나만의 일기장", "일기를 CBT, AI 등과 함께 상호작용해요", "https://deario.toy-project.n-e.kr"),
					ServiceCard("나만의 TODO 앱", "나만의 할 일을 기록해보세요", "https://development-support.p-e.kr"),
				),
			),
		},
	})
}
