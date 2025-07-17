package views

import (
	"simple-server/pkg/util/gomutil"
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
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
			Header(Class("fixed yellow4"),
				Nav(
					A(Href("/"), Class("max"),
						H3(Text("Deario")),
					),
					A(ID("login"), Href("/login"), x.Data(""), x.Show("!$store.auth.isAuthed"), Text("Login")),
					A(ID("logout"), Href("#"), x.Data(""), x.Show("$store.auth.isAuthed"), Attr("onclick", "logoutUser()"), Text("Logout")),
				),
			),
			Main(Class("responsive"),
				H5(Text("월별 일기 통계")),
				Canvas(ID("countChart")),
				H5(Text("월별 평균 기분")),
				Canvas(ID("moodChart")),
				H5(Text("월별 기분 분포")),
				Canvas(ID("moodStackChart")),
			),
		},
	})
}
