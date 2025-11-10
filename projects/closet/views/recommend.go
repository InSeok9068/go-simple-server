package views

import (
	"fmt"
	"strconv"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	h "maragu.dev/gomponents-htmx"
)

type RecommendationItem struct {
	Kind string
	Item ClosetItem
}

var recommendationKindOrder = []string{"top", "bottom", "shoes", "accessory"}

func RecommendationDialog(results []RecommendationItem, weather, style, cacheToken string, hasMore bool, locks map[string]int64) Node {
	itemMap := make(map[string]ClosetItem, len(results))
	for _, result := range results {
		itemMap[result.Kind] = result.Item
	}

	rows := make([]Node, 0, len(recommendationKindOrder))
	for _, kind := range recommendationKindOrder {
		item, ok := itemMap[kind]
		if ok {
			itemCopy := item
			rows = append(rows, renderRecommendationRow(kind, &itemCopy, locks[kind]))
			continue
		}
		rows = append(rows, renderRecommendationRow(kind, nil, 0))
	}

	body := []Node{
		H5(Text("추천 결과")),
	}

	if len(results) == 0 {
		body = append(body,
			P(Class("caption"), Text("추천 가능한 옷이 아직 없어요.")),
			Div(Class("row end"),
				Button(Class("button outline"), Type("button"),
					Attr("onclick", "this.closest('dialog').remove()"),
					Text("닫기"),
				),
			),
		)
	} else {
		body = append(body,
			Form(Class("stack gap-md"),
				h.Post("/recommend"),
				h.Target("#recommend-result"),
				h.Swap("innerHTML"),
				Input(Type("hidden"), Name("weather"), Value(weather)),
				Input(Type("hidden"), Name("style"), Value(style)),
				Input(Type("hidden"), Name("skip_ids"), Value(cacheToken)),
				Div(Class("stack gap-sm"), Group(rows)),
				Div(Class("row end gap-sm"),
					Button(Class("button outline"), Type("button"),
						Attr("onclick", "this.closest('dialog').remove()"),
						Text("닫기"),
					),
					If(hasMore,
						Button(Class("button"), Type("submit"), Text("다시 추천받기")),
					),
					If(!hasMore,
						Button(Class("button"), Type("button"), Disabled(),
							Text("더 추천할 옷이 없어요"),
						),
					),
				),
			),
		)
	}

	return Dialog(Class("active recommend-dialog"),
		DataAttr("recommend-cache", cacheToken),
		Div(Class("stack gap-md"), Group(body)),
	)
}

func renderRecommendationRow(kind string, item *ClosetItem, lockID int64) Node {
	var figure Node
	var lockControl Node
	if item != nil {
		value := strconv.FormatInt(item.ID, 10)
		figure = Figure(Class("recommend-card"),
			Img(
				Src(item.ImageURL),
				Alt(fmt.Sprintf("%s 이미지", item.KindLabel)),
				Width("140"),
				Height("140"),
				Loading("lazy"),
			),
		)
		lockControl = Label(Class("checkbox lock-toggle"),
			Input(Type("checkbox"), Name(fmt.Sprintf("lock_%s", kind)), Value(value),
				If(lockID == item.ID, Checked()),
			),
			Span(Text("고정")),
		)
	} else {
		figure = Div(Class("recommend-card recommend-card--empty"),
			Span(Class("caption"), Text("추천 없음")),
		)
		lockControl = Div()
	}

	return Div(Class("recommend-row row gap-sm align-center"),
		figure,
		lockControl,
	)
}
