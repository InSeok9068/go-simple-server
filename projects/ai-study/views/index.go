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
		Head:     shared.Heads(title),
		Body: []Node{
			b.Container(b.MaxTablet, b.MarginTop(3),
				b.Columns(
					b.Column(
						b.Box(
							b.Content(
								H3(Text("오늘 무엇을 공부하고싶은지 고민이신가요 ?")),
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
							),
						),
						P(Text("⏬ 결과 확인")),
						b.Box(e.Class("content"),
							Div(ID("result")),
						),
					),
				),
			),
		},
	})
}

/*
package views

templ Page() {
	<div class="box">
		<div class="content">
			<h3>오늘 무엇을 공부하고싶은지 고민이신가요 ?</h3>
			<p>AI와 함께 공부할 수 있도록 주제와 사이트 접속 링크를 제공해드려요</p>
			<form
				hx-post="/ai-study"
				hx-target="#result"
				x-data="{ loading: false }"
				@htmx:before-request="loading = true"
				@htmx:after-on-load="loading = false"
			>
				<div class="field has-addons">
					<div class="control is-expanded">
						<input class="input" type="text" name="input" placeholder="내용을 입력해주세요." required autofocus/>
					</div>
					<div class="control">
						<button class="button is-info">
							입력
							<img x-show="loading" class="htmx-indicator" src="/static/spinner.svg"/>
						</button>
					</div>
				</div>
			</form>
			<div class="columns">
				<div class="column is-half is-offset-one-quarter">
					<button
						class="button is-focused is-fullwidth mt-3"
						type="submit"
						hx-post="/ai-study-random"
						hx-target="#result"
						x-data="{ loading: false }"
						@htmx:before-request="loading = true"
						@htmx:after-on-load="loading = false"
					>
						❔ 랜덤 주제로 시작하기
						<img x-show="loading" class="htmx-indicator" src="/static/spinner.svg"/>
					</button>
				</div>
			</div>
		</div>
	</div>
	<p>⏬ 결과 확인 </p>
	<div class="box content mt-3">
		<div class="columns is-mobile">
			<div class="column">
				<button class="button" id="copy">글 복사</button>
			</div>
			<div class="column">
				<a href="https://chatgpt.com/" target="_blank">ChatGPT 📙</a>
			</div>
			<div class="column">
				<a href="https://gemini.google.com/app?hl=ko" target="_blank">Gemini 📒</a>
			</div>
		</div>
		<div id="result"></div>
	</div>
	<script>
		document.getElementById("copy").addEventListener('click', () => {
			navigator.clipboard.writeText(document.getElementById("result").innerText).then(() => alert("복사 되었습니다."));
		})
	</script>
}

templ Index() {
	@HomepageLayout("", Page())
}

*/
