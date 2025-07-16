package views

import (
	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/db"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Setting(userSetting db.UserSetting) Node {
	return HTML5(HTML5Props{
		Title:    "설정",
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer("설정"),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("/static/deario.js")),
			},
		),
		Body: []Node{
			shared.Snackbar(),
			/* Body */
			Main(
				Form(
					FieldSet(
						Legend(Text("설정")),
						Div(Class("field border label"),
							Input(),
						),
					),
				),
			),
			/* Body */
		},
	})
}
