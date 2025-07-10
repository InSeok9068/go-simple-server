package shared

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Login() Node {
	return HTML5(HTML5Props{
		Title:    "로그인",
		Language: "ko",
		Head: append(
			HeadsWithBeer("로그인"),
			HeadWithFirebaseLogin()...,
		),
		Body: []Node{
			Main(Class("responsive"),
				Div(ID("firebaseui-auth-container")),
				Div(ID("loader"), Text("Loading...")),
			),
		}},
	)
}
