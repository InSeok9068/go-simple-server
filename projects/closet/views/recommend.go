package views

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type RecommendationItem struct {
	Kind string
	Item ClosetItem
}

func RecommendationDialog(results []RecommendationItem) Node {
	items := make([]Node, 0, len(results))
	for _, result := range results {
		item := result.Item
		items = append(items,
			Article(Class("row gap-sm align-center"),
				Img(Src(item.ImageURL), Alt(fmt.Sprintf("%s 이미지", item.KindLabel)), Width("80"), Height("80")),
				Div(Class("stack gap-xxs"),
					Strong(Text(item.KindLabel)),
					If(item.TagLine != "",
						Span(Class("caption"), Text(item.TagLine)),
					),
				),
			),
		)
	}

	return Dialog(Class("active"),
		Div(Class("stack gap-sm"),
			H3(Text("추천 결과")),
			P(Class("caption"), Text("조건에 맞춰 상의/하의/신발 순으로 골라봤어요.")),
			If(len(items) == 0,
				P(Class("caption"), Text("추천 가능한 옷이 아직 없어요.")),
			),
			If(len(items) > 0,
				Div(Class("stack gap-sm"), Group(items)),
			),
			Div(Class("row end gap-sm"),
				Button(Class("button outline"), Attr("onclick", "this.closest('dialog').remove()"), Text("닫기")),
			),
		),
	)
}
