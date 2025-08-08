package views

import (
	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	h "maragu.dev/gomponents-htmx"

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
			shared.HeadGoogleFonts("Gamja Flower"),
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
				Div(Class("large-space")),
				Button(Class("responsive small-elevate large fill"),
					h.Post("/ai-report"),
					h.Swap("none"),
					h.On("htmx:after-on-load", "if (event.detail.successful) showInfo('리포트 생성 요청이 되었습니다.')"),
					Span(Class("bold"), Text("AI 상담 리포트")),
				),
				P(Text("※ 최근 30개의 일기내용을 참고해서 리포트가 작성됩니다.")),
			),
			/* Body */
		},
	})
}
