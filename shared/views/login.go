package shared

import (
	"simple-server/pkg/util/gomutil"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Login() Node {
	return HTML5(HTML5Props{
		Title:    "로그인",
		Language: "ko",
		Head: gomutil.MergeHeads(
			HeadsDefault("로그인"),
			HeadsCustom(),
			HeadWithFirebaseLogin(),
		),
		Body: []Node{
			Main(Class("login-page"),
				Div(ID("firebaseui-auth-container")),
				Div(ID("loader"),
					Img(Src("/shared/static/spinner.svg"), Alt("로딩")),
					Span(Text("로딩 중...")),
				),
			),
		}},
	)
}
