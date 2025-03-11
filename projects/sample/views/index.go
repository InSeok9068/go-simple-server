package views

import (
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	h "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Radio() Node {
	return Div(x.Data("{ selected: 1 }"),
		Input(Type("radio"), x.Model("selected"), Value("1")),
		Input(Type("radio"), x.Model("selected"), Value("2")),
		Button(x.On("click", "selected = 3"),
			Text("3"),
		),
		P(x.Text("selected")),
	)
}

func Radio2() Group {
	return Group{
		Input(Type("radio"), x.Model("selected"), Value("1")),
		Input(Type("radio"), x.Model("selected"), Value("2")),
		Button(x.On("click", "selected = 3"),
			Text("3"),
		),
		P(x.Text("selected")),
	}
}

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: append(
			shared.HeadsWithPicoAndTailwind(title),
			Script(Src("./static/index.js")),
		),
		Body: []Node{
			Main(
				Div(Class("flex flex-col gap-2"),
					Div(Class("flex gap-2"),
						Button(h.Get("/radio"), h.Target("#box1"),
							Text("Click Me 1"),
						),
						Button(h.Get("/radio2"), h.Target("#box2"),
							Text("Click Me 2"),
						),
					),
					Article(ID("box1"), Class("flex")),
					Article(ID("box2"), Class("flex"), x.Data("{ selected: 1 }")),
				),
			),
		},
	})
}
