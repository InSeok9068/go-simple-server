package views

import (
	"simple-server/pkg/util/gomutil"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Meta(Name("description"), Content("자산 포트폴리오")),
				Link(Rel("manifest"), Href("/manifest.json")),
			},
		),
		Body: []Node{
			shared.Snackbar(),
			Main(
				Text("Hello World"),
			),
		},
	})
}
