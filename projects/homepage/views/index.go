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
		H3(Text(name)),
		Text(desc),
		Footer(Class("text-center"),
			A(Href(url),
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
		Head:     shared.HeadsWithPicoAndTailwind(title),
		Body: []Node{
			Main(
				Nav(
					Ul(
						Li(
							H2(A(Href("/"), Text(title))),
						),
					),
				),
				ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져줘요", "https://ai-study.toy-project.n-e.kr"),
				ServiceCard("나만의 TODO 앱", "나만의 할 일을 기록해보세요", "https://development-support.p-e.kr"),
			),
		},
	})
}
