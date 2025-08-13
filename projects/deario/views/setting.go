package views

import (
	"fmt"

	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
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
			shared.HeadGoogleFonts("Gamja Flower"),
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
						// 다크모드 변경
						Nav(Text("테마")),
						Label(Class("switch icon"),
							Input(Type("checkbox"),
								x.Data(""),
								x.Bind("checked", "$store.theme.value === 'dark'"),
								x.On("change", "$store.theme.toggle()")),
							Span(
								I(Text("dark_mode")),
							),
						),
						// 테마색 변경
						Nav(Text("테마색")),
						Div(Class("field border label"),
							Button(
								I(Text("palette")),
								Span(Text("Color")),
								Input(Type("color"),
									x.Data(""),
									x.Model("$store.theme.color"),
									x.On("change", "$store.theme.setColor($event.target.value)"),
								),
							),
						),
						// 폰트 선택
						Nav(Text("폰트")),
						Div(Class("field border label"),
							Select(
								x.Data(""),
								x.Model("$store.font.value"),
								x.On("change", "$store.font.set($event.target.value)"),
								Option(Value("gamja"), Text("Gamja Flower (기본)")),
								Option(Value("humanist"), Text("Humanist")),
								Option(Value("geometric_humanist"), Text("Geometric Humanist")),
								Option(Value("classical_humanist"), Text("Classical Humanist")),
								Option(Value("neo_grotesque"), Text("Neo Grotesque")),
								Option(Value("monospace_code"), Text("Monospace Code")),
								Option(Value("industrial"), Text("Industrial")),
								Option(Value("rounded_sans"), Text("Rounded Sans")),
								Option(Value("antique"), Text("Antique")),
							),
							Label(Text("폰트")),
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
