package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	h "maragu.dev/gomponents-htmx"
)

type RecommendationItem struct {
	Kind string
	Item ClosetItem
}

func RecommendationDialog(results []RecommendationItem, weather, style, cacheToken string, hasMore bool) Node {
	items := make([]Node, 0, len(results))
	for _, result := range results {
		item := result.Item
		items = append(items,
			Figure(Class("recommend-card"),
				Img(
					Src(item.ImageURL),
					Alt(item.KindLabel+" 이미지"),
					Width("160"),
					Height("160"),
					Loading("lazy"),
				),
			),
		)
	}

	return Dialog(Class("active recommend-dialog"),
		DataAttr("recommend-cache", cacheToken),
		Div(Class("stack gap-md"),
			H3(Text("추천 결과")),
			If(len(items) == 0,
				P(Class("caption"), Text("추천 가능한 옷이 아직 없어요.")),
			),
			If(len(items) > 0,
				Div(Class("recommend-grid row gap-sm wrap center"),
					Group(items),
				),
			),
			Div(Class("row end gap-sm"),
				If(hasMore && len(items) > 0,
					Form(Class("row gap-sm"),
						h.Post("/recommend"),
						h.Target("#recommend-result"),
						h.Swap("innerHTML"),
						Input(Type("hidden"), Name("weather"), Value(weather)),
						Input(Type("hidden"), Name("style"), Value(style)),
						Input(Type("hidden"), Name("skip_ids"), Value(cacheToken)),
						Button(Class("button"), Text("다시 추천받기")),
					),
				),
				Button(Class("button outline"), Type("button"),
					Attr("onclick", "this.closest('dialog').remove()"),
					Text("닫기"),
				),
			),
		),
	)
}
