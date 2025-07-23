package views

import (
	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Statistic() Node {
	return HTML5(HTML5Props{
		Title:    "통계",
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer("통계"),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("https://cdn.jsdelivr.net/npm/chart.js")),
				Script(Src("/static/statistic.js")),
			},
		),

		Body: []Node{
			shared.Snackbar(),
			/* Header */
			layout.AppHeader(),
			/* Header */

			/* Body */
			Main(Class("responsive"),
				H5(Text("월별 일기 통계")),
				Canvas(ID("countChart")),
				H5(Text("월별 기분 분포")),
				Canvas(ID("moodStackChart")),
			),
			/* Body */
		},
	})
}
