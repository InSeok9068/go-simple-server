package views

import (
	lucide "github.com/eduardolat/gomponents-lucide"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	shared "simple-server/shared/views"
)

func ServiceCard(name string, desc string, url string) Node {
	return Article(
		H5(Text(name)),
		P(Text(desc)),
		Hr(),
		Div(Class("space")),
		Nav(
			A(Href(url), Class("underline"),
				lucide.ExternalLink(),
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
			Nav(Class("bottom s"),
				A(I(Text("home"))),
				A(I(Text("search"))),
				A(I(Text("share"))),
				A(I(Text("more_vert"))),
			),

			Header(Class("fixed responsive yellow4"),
				A(Href("/"),
					H3(Text(title)),
				),
			),

			Div(Class("space")),

			Main(Class("responsive"),
				ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져줘요", "https://ai-study.toy-project.n-e.kr"),
				ServiceCard("나만의 TODO 앱", "나만의 할 일을 기록해보세요", "https://development-support.p-e.kr"),
			),
		},
	})
}
