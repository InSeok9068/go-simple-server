package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DiaryImages(img1, img2, img3 string) Node {
	nodes := []Node{}
	urls := []string{img1, img2, img3}
	for _, u := range urls {
		if u != "" {
			nodes = append(nodes, Img(Src(u), Class("responsive")))
		}
	}
	if len(nodes) == 0 {
		nodes = append(nodes, P(Text("이미지가 없습니다.")))
	}
	all := append([]Node{ID("diary-image-content")}, nodes...)
	return Div(all...)
}
