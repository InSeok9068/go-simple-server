package views

import (
	lucide "github.com/eduardolat/gomponents-lucide"
	_ "github.com/glsubri/gomponents-alpine"
	b "github.com/willoma/bulma-gomponents"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	shared "simple-server/shared/views"
)

func ServiceCard(name string, desc string, url string) Node {
	return b.Card(b.Padding(3),
		b.CardHeader(
			b.CardHeaderTitle(Text(name)),
			b.CardHeaderIcon(),
		),
		b.Content(desc),
		b.CardFooter(
			A(Href(url), Class("card-footer-item"), Text("서비스 이동"), lucide.ArrowUpRight())))
}

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head:     shared.Heads(title),
		Body: []Node{
			b.Container(
				b.MaxTablet,
				b.Navbar(
					b.NavbarBrand(
						b.NavbarAHref("/", b.FontSize(3), Text("홈페이지")))),
				ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져줘요", "https://ai-study.toy-project.n-e.kr"),
				ServiceCard("나만의 TODO 앱", "나만의 할 일을 기록해보세요", "https://development-support.p-e.kr"),
			),
		},
	})
}
