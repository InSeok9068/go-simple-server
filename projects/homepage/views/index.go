package views

import (
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// ServiceCard는 BeerCSS 카드 형태로 서비스 정보를 표시합니다.
func ServiceCard(name, desc, url string) Node {
	return Article(Class("card"),
		H6(Text(name)),
		P(Text(desc)),
		Nav(
			A(Class("chip"), Href(url), Text("자세히 보기")),
		),
	)
}

// FeatureCard는 서비스의 특징을 한눈에 보여줍니다.
func FeatureCard(icon, title, desc string) Node {
	return Article(Class("card center"),
		I(Class("large material-icons"), Text(icon)),
		H6(Text(title)),
		P(Text(desc)),
	)
}

// Index는 홈페이지 메인 화면을 렌더링합니다.
func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: append(
			shared.HeadsWithBeer(title),
			Link(Rel("preconnect"), Href("https://fonts.googleapis.com")),
			Link(Rel("preconnect"), Href("https://fonts.gstatic.com"), CrossOrigin("anonymous")),
			Link(Rel("stylesheet"), Href("https://fonts.googleapis.com/css2?family=Noto+Sans+KR:wght@300;400;500;600;700&display=swap")),
			Link(Rel("stylesheet"), Href("https://fonts.googleapis.com/icon?family=Material+Icons")),
		),
		Body: []Node{
			Header(Class("appbar"),
				Nav(
					A(Class("brand"), Href("/"), H3(Text(title))),
					Div(Class("max")),
					A(Href("/login"), Class("chip"), Text("로그인")),
				),
			),

			Main(
				// 히어로 영역
				Section(Class("header primary center-align"),
					H1(Text("디지털 라이프를 더 편리하게")),
					P(Text("다양한 서비스를 통해 일상과 업무의 효율성을 높여보세요.")),
					Nav(
						A(Href("/services"), Class("button"), Text("서비스 살펴보기")),
						A(Href("/contact"), Class("button outline"), Text("문의하기")),
					),
				),

				// 특징 영역
				Section(Class("container"),
					H3(Class("center-align"), Text("특징")),
					Div(Class("row"),
						FeatureCard("bolt", "빠른 속도", "가벼운 구조로 신속하게 작동합니다."),
						FeatureCard("smartphone", "모바일 최적화", "모든 기기에서 깔끔하게 보입니다."),
						FeatureCard("palette", "세련된 테마", "BeerCSS를 활용해 감각적인 디자인을 제공합니다."),
					),
				),

				// 서비스 소개
				Section(Class("container"),
					H3(Class("center-align"), Text("서비스")),
					Div(Class("row"),
						ServiceCard("AI 공부 도우미", "AI가 공부 주제를 제안하고 학습을 도와주는 서비스입니다.", "https://ai-study.toy-project.n-e.kr"),
						ServiceCard("나만의 일기장", "일기를 기록하고 AI와 함께 감정을 관리해 보세요.", "https://deario.toy-project.n-e.kr"),
						ServiceCard("나만의 TODO 앱", "할 일을 효율적으로 관리할 수 있는 투두리스트입니다.", "https://development-support.p-e.kr"),
					),
				),
			),

			Footer(Class("appbar"),
				P(Class("center-align"), Text("© 2025 모든 권리 보유.")),
			),
		},
	})
}
