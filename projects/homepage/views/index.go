package views

import (
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// ServiceCard는 서비스 항목을 카드 형태로 표시합니다
func ServiceCard(name string, desc string, url string) Node {
	return Div(Class("group relative overflow-hidden rounded-xl bg-white shadow-md transition duration-300 hover:shadow-lg hover:translate-y-[-4px]"),
		A(Href(url), Class("absolute inset-0 z-10"), Aria("label", name)),
		Div(Class("p-6"),
			// 아이콘 영역
			Div(Class("mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-indigo-100 text-indigo-600 group-hover:bg-indigo-600 group-hover:text-white transition-colors duration-300"),
				Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
				</svg>`),
			),
			// 내용 영역
			H3(Class("mt-4 text-xl font-bold text-gray-900 group-hover:text-indigo-600 transition-colors duration-300"), Text(name)),
			P(Class("mt-2 text-gray-600 break-words"), Text(desc)),
			// 화살표 부분
			Div(Class("mt-6 flex items-center text-sm font-medium text-indigo-600 group-hover:text-indigo-700"),
				Span(Class("mr-2"), Text("자세히 보기")),
				Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
				</svg>`),
			),
		),
	)
}

// Index는 메인 페이지를 렌더링합니다
func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: append(shared.HeadsWithTailwind(title),
			Link(Rel("preconnect"), Href("https://fonts.googleapis.com")),
			Link(Rel("preconnect"), Href("https://fonts.gstatic.com"), CrossOrigin("anonymous")),
			Link(Rel("stylesheet"), Href("https://fonts.googleapis.com/css2?family=Noto+Sans+KR:wght@300;400;500;600;700&display=swap")),
		),
		Body: []Node{
			Div(Class("font-sans bg-gray-50 text-gray-900 min-h-screen flex flex-col overflow-x-hidden"),
				x.Data(`{ isMenuOpen: false, isScrolled: false }`),

				// 네비게이션 바 - 스크롤 시 배경 변경
				Nav(Class("fixed top-0 left-0 right-0 z-50 transition-all duration-300 w-full"),
					x.Class(`{ 'bg-white shadow-md': isScrolled, 'bg-transparent': !isScrolled }`),
					x.On("scroll.window", `isScrolled = window.scrollY > 20`),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("flex h-16 items-center justify-between"),
							// 로고 부분
							A(Href("/"), Class("flex items-center"),
								Div(Class("flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md bg-indigo-600 text-white mr-3"),
									Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
									</svg>`),
								),
								Span(Class("text-lg font-bold truncate"), Text(title)),
							),

							// 데스크톱 메뉴
							Div(Class("hidden md:flex items-center space-x-6"),
								A(Href("/"), Class("font-medium text-gray-900 hover:text-indigo-600 transition-colors"), Text("홈")),
								A(Href("/about"), Class("font-medium text-gray-900 hover:text-indigo-600 transition-colors"), Text("소개")),
								A(Href("/services"), Class("font-medium text-gray-900 hover:text-indigo-600 transition-colors"), Text("서비스")),
								A(Href("/contact"), Class("font-medium text-gray-900 hover:text-indigo-600 transition-colors"), Text("문의하기")),
								A(Href("/login"), Class("ml-4 inline-flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"), Text("로그인")),
							),

							// 모바일 메뉴 버튼
							Button(Type("button"),
								Class("inline-flex items-center justify-center rounded-md p-2 text-gray-900 md:hidden"),
								Aria("expanded", "false"),
								x.On("click", "isMenuOpen = !isMenuOpen"),
								Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
								</svg>`),
							),
						),
					),
				),

				// 모바일 메뉴
				Div(Class("fixed inset-0 z-40 transform transition-all duration-300 md:hidden"),
					x.Class(`{ 'translate-x-0 ease-out': isMenuOpen, '-translate-x-full ease-in': !isMenuOpen }`),
					// 배경 오버레이
					Div(Class("absolute inset-0 bg-black bg-opacity-50 transition-opacity"),
						x.Class(`{ 'opacity-100': isMenuOpen, 'opacity-0': !isMenuOpen }`),
						x.On("click", "isMenuOpen = false"),
					),
					// 사이드바 메뉴
					Div(Class("relative max-w-xs w-full bg-white h-full shadow-xl flex flex-col overflow-y-auto"),
						Div(Class("flex items-center justify-between p-4 border-b"),
							Span(Class("text-lg font-medium"), Text("메뉴")),
							Button(Type("button"),
								Class("rounded-md p-2 text-gray-400 hover:text-gray-500"),
								x.On("click", "isMenuOpen = false"),
								Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>`),
							),
						),
						// 메뉴 항목들
						Div(Class("flex-1 px-2 py-4 space-y-1"),
							A(Href("/"), Class("block rounded-md px-3 py-2 text-base font-medium text-gray-900 hover:bg-gray-100 hover:text-indigo-600"), Text("홈")),
							A(Href("/about"), Class("block rounded-md px-3 py-2 text-base font-medium text-gray-900 hover:bg-gray-100 hover:text-indigo-600"), Text("소개")),
							A(Href("/services"), Class("block rounded-md px-3 py-2 text-base font-medium text-gray-900 hover:bg-gray-100 hover:text-indigo-600"), Text("서비스")),
							A(Href("/contact"), Class("block rounded-md px-3 py-2 text-base font-medium text-gray-900 hover:bg-gray-100 hover:text-indigo-600"), Text("문의하기")),
						),
						// 로그인 버튼
						Div(Class("px-4 py-4 border-t"),
							A(Href("/login"), Class("flex w-full items-center justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700"), Text("로그인")),
						),
					),
				),

				// 히어로 섹션
				Section(Class("relative pt-16 pb-20 md:pt-24 md:pb-28 bg-gradient-to-br from-indigo-900 to-purple-800 overflow-hidden"),
					// 배경 패턴
					Div(Class("absolute inset-0 opacity-10"),
						Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="100%" height="100%">
							<defs>
								<pattern id="grid-pattern" width="40" height="40" patternUnits="userSpaceOnUse">
									<path d="M0 40L40 0M20 40L40 20M0 20L20 0" stroke="white" stroke-width="1" fill="none" />
								</pattern>
							</defs>
							<rect width="100%" height="100%" fill="url(#grid-pattern)" />
						</svg>`),
					),
					Div(Class("relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 text-center"),
						Div(Class("mx-auto max-w-3xl text-center"),
							H1(Class("text-4xl font-bold tracking-tight text-white sm:text-5xl lg:text-6xl"),
								Text("디지털 라이프를 더 편리하게"),
							),
							P(Class("mx-auto mt-6 max-w-xl text-lg text-indigo-100"),
								Text("다양한 서비스를 통해 일상과 업무의 효율성을 높여보세요. 저희 플랫폼에서 제공하는 다양한 도구로 생산성을 높이세요."),
							),
							Div(Class("mt-10 flex flex-col sm:flex-row gap-4 justify-center"),
								A(Href("/services"),
									Class("rounded-md bg-white px-6 py-3 text-base font-medium text-indigo-700 shadow-sm hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600"),
									Text("서비스 살펴보기"),
								),
								A(Href("/contact"),
									Class("rounded-md border border-transparent border-white bg-transparent px-6 py-3 text-base font-medium text-white hover:bg-white/10 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600"),
									Text("문의하기"),
								),
							),
						),
					),
				),

				// 서비스 섹션
				Section(Class("py-16 md:py-24 bg-white"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("mx-auto max-w-3xl text-center mb-16"),
							H2(Class("text-base font-semibold uppercase tracking-wide text-indigo-600"), Text("서비스")),
							P(Class("mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"), Text("당신을 위한 맞춤 서비스")),
							P(Class("mt-4 text-lg text-gray-500"), Text("일상과 업무에 필요한 다양한 서비스를 제공합니다.")),
						),

						// 서비스 카드 그리드
						Div(Class("grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3"),
							ServiceCard("AI 공부 도우미", "AI가 공부 주제를 던져주고 학습을 도와주는 서비스입니다. 다양한 주제로 학습해보세요.", "https://ai-study.toy-project.n-e.kr"),
							ServiceCard("나만의 일기장", "일기를 CBT(인지행동치료), AI 등과 함께 상호작용하며 감정을 관리하는 서비스입니다.", "https://deario.toy-project.n-e.kr"),
							ServiceCard("나만의 TODO 앱", "할 일을 효율적으로 관리하고 완료하는 기쁨을 느낄 수 있는 투두리스트입니다.", "https://development-support.p-e.kr"),
						),
					),
				),

				// 특징 섹션
				Section(Class("py-16 md:py-24 bg-gray-50"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("mx-auto max-w-3xl text-center mb-16"),
							H2(Class("text-base font-semibold uppercase tracking-wide text-indigo-600"), Text("특징")),
							P(Class("mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"), Text("왜 저희 서비스가 특별할까요?")),
							P(Class("mt-4 text-lg text-gray-500"), Text("직관적인 디자인부터 강력한 기능까지, 다양한 이점을 제공합니다.")),
						),

						// 특징 그리드
						Div(Class("grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3"),
							// 특징 1
							Div(Class("relative rounded-2xl bg-white p-6 shadow-md"),
								Div(Class("flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white"),
									Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
									</svg>`),
								),
								H3(Class("mt-4 text-lg font-medium text-gray-900"), Text("뛰어난 성능")),
								P(Class("mt-2 text-base text-gray-500"), Text("최적화된 서비스로 빠른 응답 속도와 안정적인 성능을 제공합니다.")),
							),
							// 특징 2
							Div(Class("relative rounded-2xl bg-white p-6 shadow-md"),
								Div(Class("flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white"),
									Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
									</svg>`),
								),
								H3(Class("mt-4 text-lg font-medium text-gray-900"), Text("안전한 보안")),
								P(Class("mt-2 text-base text-gray-500"), Text("최신 보안 기술로 사용자의 데이터를 안전하게 보호합니다.")),
							),
							// 특징 3
							Div(Class("relative rounded-2xl bg-white p-6 shadow-md"),
								Div(Class("flex h-12 w-12 items-center justify-center rounded-md bg-indigo-500 text-white"),
									Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>`),
								),
								H3(Class("mt-4 text-lg font-medium text-gray-900"), Text("24시간 지원")),
								P(Class("mt-2 text-base text-gray-500"), Text("언제든지 도움이 필요할 때 지원팀이 함께합니다.")),
							),
						),
					),
				),

				// CTA 섹션
				Section(Class("py-16 md:py-24 bg-indigo-700"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("mx-auto max-w-3xl text-center"),
							H2(Class("text-3xl font-bold tracking-tight text-white sm:text-4xl"), Text("지금 바로 시작하세요")),
							P(Class("mt-4 text-lg text-indigo-100"), Text("다양한 서비스를 무료로 이용해보고 더 많은 기능을 경험해보세요.")),
							Div(Class("mt-8 flex justify-center"),
								A(Href("/signup"), Class("rounded-md bg-white px-6 py-3 text-base font-medium text-indigo-700 shadow-sm hover:bg-indigo-50"), Text("무료로 시작하기")),
							),
						),
					),
				),

				// 푸터
				Footer(Class("bg-gray-900 py-12"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4"),
							// 회사 정보
							Div(Class(""),
								A(Href("/"), Class("flex items-center"),
									Div(Class("flex h-8 w-8 items-center justify-center rounded-md bg-indigo-600 text-white mr-3"),
										Raw(`<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
											<path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
										</svg>`),
									),
									Span(Class("text-lg font-bold text-white truncate"), Text(title)),
								),
								P(Class("mt-4 text-base text-gray-400"), Text("더 나은 디지털 라이프를 위한 다양한 서비스를 제공합니다.")),
								Div(Class("mt-6 flex space-x-6"),
									A(Href("#"), Class("text-gray-400 hover:text-white"),
										Raw(`<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
											<path fill-rule="evenodd" d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z" clip-rule="evenodd" />
										</svg>`),
									),
									A(Href("#"), Class("text-gray-400 hover:text-white"),
										Raw(`<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
											<path d="M8.29 20.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0022 5.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.072 4.072 0 012.8 9.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 012 18.407a11.616 11.616 0 006.29 1.84" />
										</svg>`),
									),
									A(Href("#"), Class("text-gray-400 hover:text-white"),
										Raw(`<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
											<path fill-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clip-rule="evenodd" />
										</svg>`),
									),
								),
							),

							// 링크 그룹 1
							Div(Class(""),
								H3(Class("text-base font-medium text-white"), Text("서비스")),
								Ul(Class("mt-4 space-y-2"),
									Li(Class(""),
										A(Href("/services/ai-study"), Class("text-base text-gray-400 hover:text-white"), Text("AI 공부 도우미")),
									),
									Li(Class(""),
										A(Href("/services/diary"), Class("text-base text-gray-400 hover:text-white"), Text("나만의 일기장")),
									),
									Li(Class(""),
										A(Href("/services/todo"), Class("text-base text-gray-400 hover:text-white"), Text("나만의 TODO 앱")),
									),
								),
							),

							// 링크 그룹 2
							Div(Class(""),
								H3(Class("text-base font-medium text-white"), Text("회사 정보")),
								Ul(Class("mt-4 space-y-2"),
									Li(Class(""),
										A(Href("/about"), Class("text-base text-gray-400 hover:text-white"), Text("소개")),
									),
									Li(Class(""),
										A(Href("/team"), Class("text-base text-gray-400 hover:text-white"), Text("팀")),
									),
									Li(Class(""),
										A(Href("/blog"), Class("text-base text-gray-400 hover:text-white"), Text("블로그")),
									),
								),
							),

							// 링크 그룹 3
							Div(Class(""),
								H3(Class("text-base font-medium text-white"), Text("지원")),
								Ul(Class("mt-4 space-y-2"),
									Li(Class(""),
										A(Href("/contact"), Class("text-base text-gray-400 hover:text-white"), Text("문의하기")),
									),
									Li(Class(""),
										A(Href("/faq"), Class("text-base text-gray-400 hover:text-white"), Text("자주 묻는 질문")),
									),
									Li(Class(""),
										A(Href("/privacy"), Class("text-base text-gray-400 hover:text-white"), Text("개인정보 처리방침")),
									),
								),
							),
						),

						// 저작권 섹션
						Div(Class("mt-12 border-t border-gray-800 pt-8"),
							P(Class("text-base text-gray-400 text-center"), Text("© 2025 모든 권리 보유. 당신의 디지털 라이프를 더 편리하게.")),
						),
					),
				),
			),
		},
	})
}
