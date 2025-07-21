package views

import (
	"fmt"

	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	h "maragu.dev/gomponents-htmx"

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

			/* Header */
			layout.AppHeader(),
			/* Header */

			/* Body */
			Main(Class("responsive"),
				Form(h.Post("/setting"), h.Swap("none"),
					h.On("htmx:after-on-load", "showInfo('저장 되었습니다.')"),
					FieldSet(
						Legend(Text("설정")),
						// 알람 여부
						Nav(Text("알림")),
						Label(Class("switch"),
							Input(Type("checkbox"), Name("is_push"), Value("1"),
								If(userSetting.IsPush == 1, Checked())),
							Span(Text("")),
						),
						// 알림 시간
						Div(Class("field border label"),
							Input(Type("time"), Name("push_time"), Value(userSetting.PushTime)),
							Label(Text("알림시간")),
						),
						// 랜덤 일자 범위
						Div(Class("field border label"),
							Input(Type("number"), Name("random_range"),
								Value(fmt.Sprintf("%d", userSetting.RandomRange))),
							Label(Text("랜덤일자")),
						),
						// 핀 잠금
						Nav(Text("핀 잠금")),
						Label(Class("switch"),
							Input(Type("checkbox"), Name("pin_enabled"), Value("1"),
								If(userSetting.PinEnabled == 1, Checked())),
							Span(Text("")),
						),
						Div(Class("field border label"),
							Input(Type("password"), Name("pin"), Placeholder("핀번호")),
							Label(Text("핀번호")),
						),
						Div(Class("field border label"),
							Input(Type("number"), Name("pin_cycle"),
								Value(fmt.Sprintf("%d", userSetting.PinCycle))),
							Label(Text("핀 주기(-1:없음,0:매번)")),
						),
						Nav(Class("right-align"),
							Button(Text("저장")),
						),
					),
				),
			),
			/* Body */
		},
	})
}
