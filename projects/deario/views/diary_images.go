package views

import (
	"fmt"

	. "maragu.dev/gomponents"
	h "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func DiaryImages(date, img1, img2, img3 string) Node {
	nodes := []Node{}
	urls := []string{img1, img2, img3}
	for i, u := range urls {
		if u != "" {
			nodes = append(nodes,
				Div(
					Style("position:relative; display:inline-block;"),
					Img(Src(u), Class("responsive")),
					Button(
						Class("transparent"),
						Type("button"),
						Attr("style", "position:absolute; right:1px;"),
						h.Delete("/diary/image"),
						h.Target("#diary-image-content"),
						h.Swap("outerHTML"),
						Attr("hx-vals", fmt.Sprintf("{\"date\":\"%s\",\"slot\":%d}", date, i+1)),
						I(Text("close")),
					),
				),
				Div(Class("space")),
			)
		}
	}
	if len(nodes) == 0 {
		nodes = append(nodes, P(Text("이미지가 없습니다.")))
	}
	all := append([]Node{ID("diary-image-content")}, nodes...)
	return Div(all...)
}
