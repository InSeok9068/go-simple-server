package views

import (
	"simple-server/pkg/util/gomutil"
	"simple-server/projects/deario/views/layout"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// Privacy는 개인정보 처리방침 페이지를 렌더링한다.
func Privacy() Node {
	return HTML5(HTML5Props{
		Title:    "개인정보 처리방침",
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer("개인정보 처리방침"),
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
				H5(Text("개인정보 처리방침")),
				P(Text("Deario는 사용자의 개인정보를 소중히 여기며 다음과 같이 처리합니다.")),
				Ul(Class("list border"),
					Li(Text("- 수집 항목: 이메일, 닉네임, 일기 내용")),
					Li(Text("- 이용 목적: 서비스 제공, 맞춤형 피드백 및 알림")),
					Li(Text("- 보관 기간: 회원 탈퇴 시 즉시 파기")),
					Li(Text("- 제3자 제공: 제공하지 않음")),
					Li(Text("- 문의: dlstjr9068@gmail.com")),
				),
			),
			/* Body */
		},
	})
}
