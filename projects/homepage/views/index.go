package views

import (
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func ServiceCard(name string, desc string, url string) Node {
	return Div(Class("card large padding"),
		Div(Class("row"),
			Div(Class("max"),
				H5(Class("primary-text"), Text(name)),
				P(Class("secondary-text"), Text(desc)),
			),
			A(Href(url), Class("circle primary"),
				I(Text("arrow_forward")),
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
			Div(x.Data(`{ open: false }`),
				// 상단 앱 바
				Header(Class("fixed"),
					Nav(Class("max"),
						A(Href("/"), Class("button circle transparent"),
							I(Text("home")),
						),
						Div(Class("max")),
						A(Class("button circle transparent"),
							I(Text("search")),
						),
						A(x.On("click", "open = !open"), Class("button circle transparent"),
							I(Text("more_vert")),
						),
					),
				),

				// 사이드 메뉴
				Dialog(ID("dialog-right"), Class("right"),
					x.Class("{ 'active': open }"),
					Nav(Class("column"),
						H4(Class("padding"), Text("메뉴")),
						A(Href("/"), Class("button border"),
							I(Text("home")), Text("홈"),
						),
						A(Href("/about"), Class("button border"),
							I(Text("info")), Text("소개"),
						),
						A(Href("/contact"), Class("button border"),
							I(Text("email")), Text("문의하기"),
						),
					),
				),

				// 메인 콘텐츠
				Main(Class("responsive padding"),
					Div(Class("section"),
						H4(Class("primary-text"), Text("서비스 목록")),
						P(Text("제공 중인 서비스를 확인해보세요")),
						Div(Class("space")),
						Div(Class("grid"),
							Div(Class("s12 m6 l4"),
								ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져줘요", "https://ai-study.toy-project.n-e.kr"),
							),
							Div(Class("s12 m6 l4"),
								ServiceCard("나만의 일기장", "일기를 CBT, AI 등과 함께 상호작용해요", "https://deario.toy-project.n-e.kr"),
							),
							Div(Class("s12 m6 l4"),
								ServiceCard("나만의 TODO 앱", "나만의 할 일을 기록해보세요", "https://development-support.p-e.kr"),
							),
						),
					),
				),

				// 하단 푸터
				Footer(Class("padding"),
					Div(Class("center-align"),
						P(Class("small-text"), Text(" 2025 . All rights reserved.")),
					),
				),
			),
		},
	})
}
