package views

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	h "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	shared "simple-server/shared/views"
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
		Head:     shared.HeadsWithPicoAndTailwind(title),
		Body: []Node{
			Div(Class("container"),
				Button(h.Get("/radio"), h.Target("#box1"),
					Text("Click Me 1"),
				),
				Button(h.Get("/radio2"), h.Target("#box2"),
					Text("Click Me 2"),
				),
				Article(ID("box1")),
				Article(ID("box2"), x.Data("{ selected: 1 }")),
			),
		},
	})
}
