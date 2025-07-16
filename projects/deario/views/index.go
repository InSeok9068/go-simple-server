package views

import (
	"fmt"
	"simple-server/pkg/util/dateutil"
	"simple-server/pkg/util/gomutil"
	shared "simple-server/shared/views"
	"time"

	x "github.com/glsubri/gomponents-alpine"
	h "maragu.dev/gomponents-htmx"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string, date string, mood string) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("/static/deario.js")),
			},
		),
		Body: []Node{
			shared.Snackbar(),
			/* Header */
			Header(Class("fixed yellow4"),
				Nav(
					A(Href("/"), Class("max"),
						H3(Text("Deario")),
					),
					A(ID("login"), Href("/login"), x.Data(""),
						x.Show("!$store.auth.isAuthed"),
						Text("Login"),
					),
					A(ID("logout"), Href("#"), x.Data(""),
						x.Show("$store.auth.isAuthed"),
						Attr("onclick", "logoutUser()"),
						Text("Logout"),
					),
				),
			),
			/* Header */

			/* Body */
			Main(
				Nav(
					A(Href(fmt.Sprintf("/?date=%s", dateutil.MustAddDaysToDate(date, -1))),
						I(Text("arrow_back_ios")),
					),
					A(Href(fmt.Sprintf("/?date=%s", time.Now().Format("20060102"))),
						Button(Class("small-round small"), Text("Today")),
					),
					A(Href(fmt.Sprintf("/?date=%s", dateutil.MustAddDaysToDate(date, 1))),
						I(Text("arrow_forward_ios")),
					),
					Div(Class("max"),
						Text(DateView(date)),
					),
					I(x.Data(""), x.Text("$store.save.isOk ? 'check' : 'sync'")),
					P(Class("bold ellipsis"),
						x.Data(""),
						x.Text("$store.auth?.user?.displayName??'게스트'"),
					),
				),
				Hr(Class("medium")),

				// 일기장 조회
				Form(h.Get("/diary"), h.Target("#diary"), h.Trigger("load"), h.Swap("outerHTML"),
					Input(Type("hidden"), Name("date"), Value(date)),
				),

				DiaryContentForm(date, ""),

				Nav(
					P(Class("bold"),
						Text("오늘의 기분을 아이콘으로 표현해주세요 ✅"),
					),
				),

				Nav(
					Div(x.Data("{ mood: '"+mood+"' }"),
						Form(h.Post("/diary/mood"),
							h.Swap("none"),
							h.On("htmx:after-on-load", "if (event.detail.successful) showInfo('저장 되었습니다.')"),
							Input(Type("hidden"), Name("date"), Value(date)),
							Input(Type("hidden"), Name("mood"), x.Model("mood")),
							Button(Class("chip circle"), x.Class("mood === '1' ? 'primary' : ''"),
								x.On("click", "mood = '1'"),
								I(Text("sentiment_very_satisfied")),
							),
							Button(Class("chip circle"), x.Class("mood === '2' ? 'primary' : ''"),
								x.On("click", "mood = '2'"),
								I(Text("sentiment_satisfied")),
							),
							Button(Class("chip circle"), x.Class("mood === '3' ? 'primary' : ''"),
								x.On("click", "mood = '3'"),
								I(Text("sentiment_neutral")),
							),
							Button(Class("chip circle"), x.Class("mood === '4' ? 'primary' : ''"),
								x.On("click", "mood = '4'"),
								I(Text("sentiment_frustrated")),
							),
							Button(Class("chip circle"), x.Class("mood === '5' ? 'primary' : ''"),
								x.On("click", "mood = '5'"),
								I(Text("sentiment_extremely_dissatisfied")),
							),
						),
					),
				),

				Nav(
					Button(Class("chip circle"),
						I(Text("psychology_alt")),
						Div(Class("tooltip right"),
							Ul(Class("list` no-space"),
								Li(Text("- 짧게라도 하루를 되돌아보기")),
								Li(Text("- 감정을 솔직하게 작성")),
							),
						),
					),
					Button(Class("chip circle"), Attr("data-ui", "#cbt-dialog"),
						Text("CBT"),
					),
					Div(Class("max")),
					Img(ID("feedback-loading"), Class("htmx-indicator"), Src("/shared/static/spinner.svg")),
					Button(Attr("data-ui", "#ai-feedback"),
						Span(
							Text("일기 요정"),
						),
						Menu(Class("top no-wrap"), ID("ai-feedback"), Attr("data-ui", "#ai-feedback"),
							aiFeedbackButton("1", "칭찬받기"),
							aiFeedbackButton("2", "위로받기"),
							aiFeedbackButton("3", "충고받기"),
							aiFeedbackButton("4", "그림일기"),
							Hr(),
							Li(
								h.Get("/ai-feedback?date="+date),
								h.Target("#ai-feedback-content"),
								h.On("htmx:after-on-load", "if (event.detail.successful) showModal('#ai-feedback-dialog')"),
								Text("다시보기"),
							),
						),
					),
				),
			),
			/* Body */

			/* Footer */
			Nav(Class("bottom"),
				A(
					I(Text("calendar_month")),
					Input(Type("date"), Name("date"), Attr("onchange", "location.href='/?date=' + this.value")),
					Text("달력"),
				),
				A(Attr("data-ui", "#diary-list-dialog"),
					I(Text("list_alt")),
					Text("목록"),
				),
				A(h.Get("/diary/random"),
					I(Text("casino")),
					Text("랜덤"),
				),
				A(Attr("data-ui", "#menu-dialog"),
					I(Text("menu_open")),
					Text("메뉴"),
				),
			),
			/* Footer */

			/* CBT 인지행동치료 Dialog */
			Dialog(ID("cbt-dialog"), Class("max"),
				H5(Text("CBT 인지행동치료")),
				Nav(
					A(Class("link"), Href("https://gemini.google.com/gem/3ddb44f68f79"),
						Text("Gemini CBT Gem으로 이동하기")),
				),
				P(Class("bold"), Text("CBT 일기")),
				Ul(Class("list no-space"),
					Li(Text("떠오른 생각: 오늘 나를 힘들게 한 자동적 사고를 적기")),
					Li(P(Class("small-text"), Text("- 예) 발표 망쳤어 난 바보야"))),
					Li(Text("생각 검토: 그 생각이 진짜인지, 다른 관점은 없는지 묻기")),
					Li(P(Class("small-text"), Text("- 예) 정말 다 망쳤나? 잘한 건 없어?"))),
					Li(Text("새로운 생각: 더 나은 방식으로 다르게 생각해보기")),
					Li(P(Class("small-text"), Text("- 예) 아쉬운 점도 있었지만, 다음엔 더 잘할거야"))),
				),
				Hr(),
				Nav(),
				P(Class("bold"), Text("대화 자기 점검")),
				Ul(Class("list no-space"),
					Li(Text("감정 인정: 상대 감정을 충분히 받아줬는가")),
					Li(P(Class("small-text"), Text("- 예) 많이 힘들었겠네"))),
					Li(Text("숨은 의도: 말 뒤의 감정/의도를 추측했는가")),
					Li(P(Class("small-text"), Text("- 예) 혹시 피곤해서 예민한 건가?"))),
					Li(Text("판단/변명: 내 스스로 판단하거나 변명하지 않았는가")),
					Li(P(Class("small-text"), Text("- 예) 네가 서운했겠네. 미안해 (변명 X)"))),
				),
				Nav(Class("right-align"),
					Button(Attr("data-ui", "#cbt-dialog"),
						Text("확인"),
					),
				),
			),

			/* 일기요청 피드백 Dialog */
			Dialog(ID("ai-feedback-dialog"), Class("max"),
				H5(Text("일기 요정")),
				Form(
					Input(Type("hidden"), Name("date"), Value(date)),
					Div(ID("ai-feedback-content"),
						Text("안녕"),
					),
					Nav(Class("right-align"),
						Button(Class("border"),
							h.Post("/ai-feedback/save?"),
							h.Swap("none"),
							h.On("htmx:after-on-load", "if (event.detail.successful) closeModal('#ai-feedback-dialog'); showInfo('저장 되었습니다.')"),
							I(Text("save")),
							Span(Text("저장")),
						),
						Button(Type("button"), Attr("onclick", "closeModal('#ai-feedback-dialog')"),
							Text("확인"),
						),
					),
				),
			),

			/* 작성 일지 Dialog */
			Dialog(ID("diary-list-dialog"), Class("max"),
				H5(Text("작성 일지")),
				Ul(ID("diary-list-content"), Class("list border")),
				Div(x.Data("{ page : 1 }"),
					Nav(Class("center-align"),
						Button(
							I(Text("arrow_drop_down")),
							h.Get("/diary/list"),
							h.Target("#diary-list-content"),
							h.Trigger("load delay:1s, click"),
							h.Swap("beforeend"),
							Attr("@click", "page++"),
							Attr(":hx-vals", "JSON.stringify({ page: page })"),
						),
					),
					Nav(Class("right-align"),
						Button(
							Attr("data-ui", "#diary-list-dialog"),
							Text("확인"),
						),
					),
				),
			),

			/* 메뉴 Dialog */
			Dialog(ID("menu-dialog"), Class("right"),
				H5(Text("메뉴")),
				Ul(Class("list"),
					Li(
						I(Text("settings")),
						A(Class("link"), Href("/setting"), Text("설정")),
					),
				),
				Nav(Class("right-align"),
					Button(
						Attr("data-ui", "#menu-dialog"),
						Text("확인"),
					),
				),
			),
		},
	})
}

func DateView(date string) string {
	parsed, _ := time.Parse("20060102", date)
	dateStr := parsed.Format("1월 2일")
	return dateStr
}

func DiaryContentForm(date string, content string) Node {
	return Form(ID("diary"),
		Input(Type("hidden"), Name("date"), Value(date)),
		Div(Class("field textarea border deario-fit-height"),
			Textarea(
				Name("content"),
				AutoFocus(),
				x.Data(""),
				h.Post("/diary/save"),
				h.Swap("none"),
				h.Trigger("input delay:0.5s"),
				Attr("@input", "$store.save.unok()"),
				Attr("@htmx:after-request", "$store.save.ok()"),
				Text(content),
			),
		),
	)
}

func aiFeedbackButton(strType string, title string) Node {
	return Li(
		h.Post(fmt.Sprintf("/ai-feedback?type=%s", strType)),
		h.Include("[name='content']"), h.Target("#ai-feedback-content"),
		h.Indicator("#feedback-loading"),
		h.On("htmx:after-on-load", "if (event.detail.successful) showModal('#ai-feedback-dialog')"),
		Text(title),
	)
}
