package views

import (
	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Pin() Node {
	return HTML5(HTML5Props{
		Title:    "핀 입력",
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer("핀 입력"),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("/static/deario.js")),
			},
		),
		Body: []Node{
			shared.Snackbar(),
			layout.AppHeader(),
			Main(Class("responsive"),
				Form(Action("/pin"), Method("POST"),
					FieldSet(
						Legend(Text("핀 번호 입력")),
						Div(Class("field border label"),
							Input(Type("password"), Name("pin"), AutoFocus()),
							Label(Text("핀번호")),
						),
						Nav(Class("right-align"),
							Button(Text("확인")),
						),
					),
				),
			),
		},
	})
}
