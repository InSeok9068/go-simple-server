package views

import (
	x "github.com/glsubri/gomponents-alpine"
	b "github.com/willoma/bulma-gomponents"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Radio() Node {
	return Div(x.Data("{ selected: 1 }"),
		b.Radio(Name("radio"), x.Model("selected"), Value("1")),
		b.Radio(Name("radio"), x.Model("selected"), Value("2")),
		b.Button(x.On("click", "selected = 3"), Text("3")),
		P(x.Text("selected")),
	)
}

func Radio2() Group {
	return Group{
		b.Radio(Name("radio"), x.Model("selected"), Value("1")),
		b.Radio(Name("radio"), x.Model("selected"), Value("2")),
		b.Button(x.On("click", "selected = 3"), Text("3")),
		P(x.Text("selected")),
	}
}

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		//Head:     shared.Heads(title),
		Head: []Node{
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css")),
			Link(Rel("stylesheet"), Href("/shared/static/tailwindcss.css")),
		},
		Body: []Node{
			Div(Class("container"),
				Input(Type("text"))),
		},
	})
}
