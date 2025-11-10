package views

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"simple-server/projects/closet/db"

	x "github.com/glsubri/gomponents-alpine"
	h "maragu.dev/gomponents-htmx"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var kindLabels = map[string]string{
	"top":       "상의",
	"bottom":    "하의",
	"shoes":     "신발",
	"accessory": "액세서리",
}

// ClosetItem은 화면에서 표시할 아이템 정보를 담는다.
type ClosetItem struct {
	ID           int64
	Kind         string
	KindLabel    string
	Tags         []string
	TagLine      string
	Filename     string
	CreatedAt    time.Time
	CreatedLabel string
	Dimension    string
	ImageURL     string
}

type ClosetItemDetail struct {
	ClosetItem
	MetaSummary     string
	MetaSeason      string
	MetaStyle       string
	MetaColors      []string
	MetaColorsValue string
	TagsValue       string
}

// NewClosetItem은 DB 조회 결과를 뷰 모델로 변환한다.
func NewClosetItem(row db.ListItemsRow) ClosetItem {
	tags := extractTags(row.Tags)
	created := time.Unix(row.CreatedAt, 0)

	return ClosetItem{
		ID:           row.ID,
		Kind:         strings.ToLower(row.Kind),
		KindLabel:    KindLabel(row.Kind),
		Tags:         tags,
		TagLine:      formatTagLine(tags),
		Filename:     row.Filename,
		CreatedAt:    created,
		CreatedLabel: created.Local().Format("2006.01.02 15:04"),
		Dimension:    formatDimension(row.Width, row.Height),
		ImageURL:     fmt.Sprintf("/items/%d/image", row.ID),
	}
}

// NewClosetItemDetail은 상세 보기용 데이터를 구성한다.
func NewClosetItemDetail(row db.GetItemDetailRow) ClosetItemDetail {
	tags := extractTags(row.Tags)
	created := time.Unix(row.CreatedAt, 0)

	base := ClosetItem{
		ID:           row.ID,
		Kind:         strings.ToLower(row.Kind),
		KindLabel:    KindLabel(row.Kind),
		Tags:         tags,
		TagLine:      formatTagLine(tags),
		Filename:     row.Filename,
		CreatedAt:    created,
		CreatedLabel: created.Local().Format("2006.01.02 15:04"),
		Dimension:    formatDimension(row.Width, row.Height),
		ImageURL:     fmt.Sprintf("/items/%d/image", row.ID),
	}

	colors := extractTags(nullStringValue(row.MetaColors))

	return ClosetItemDetail{
		ClosetItem:      base,
		MetaSummary:     nullStringValue(row.MetaSummary),
		MetaSeason:      nullStringValue(row.MetaSeason),
		MetaStyle:       nullStringValue(row.MetaStyle),
		MetaColors:      colors,
		MetaColorsValue: strings.Join(colors, ", "),
		TagsValue:       strings.Join(tags, ", "),
	}
}

// KindLabel은 kind 코드에 대응하는 한글 라벨을 반환한다.
func KindLabel(kind string) string {
	label, ok := kindLabels[strings.ToLower(kind)]
	if ok {
		return label
	}
	return "기타"
}

// ItemsSection은 아이템 그룹 목록을 렌더링한다.
func ItemsSection(groups map[string][]ClosetItem) Node {
	sections := make([]Node, 0, len(groups))
	order := []string{"top", "bottom", "shoes", "accessory"}
	for _, kind := range order {
		items, ok := groups[kind]
		if !ok {
			continue
		}
		sections = append(sections, itemGroup(kind, items))
	}
	return Div(ID("items-list"), Group(sections))
}

func itemGroup(kind string, items []ClosetItem) Node {
	title := KindLabel(kind)
	total := len(items)

	cardNodes := make([]Node, 0, total)
	for _, item := range items {
		cardNodes = append(cardNodes, itemCard(item))
	}

	return Section(
		Nav(
			H6(Text(title)),
			Span(Class("caption"), Text(fmt.Sprintf("%d개", total))),
		),
		If(total == 0,
			Div(Class("surface-container"),
				P(Class("caption"), Text("조건에 맞는 아이템이 아직 없어요.")),
			),
		),
		If(total > 0,
			Div(Class("row scroll"), Group(cardNodes)),
		),
	)
}

func itemCard(item ClosetItem) Node {
	return Article(Class("border small-padding"),
		Attr("role", "button"),
		Attr("tabindex", "0"),
		Attr("aria-label", fmt.Sprintf("%s 상세 보기", item.KindLabel)),
		h.Get(fmt.Sprintf("/items/%d/detail", item.ID)),
		h.Target("#item-detail-content"),
		h.Swap("innerHTML"),
		h.Trigger("click, keyup[key=='Enter']"),
		Attr("hx-indicator", "#item-detail-loading"),
		Attr("hx-on::before-request", "showModal('#item-detail-dialog', true);"),
		Attr("hx-on::after-request", "if(!event.detail.successful){ closeModal('#item-detail-dialog'); }"),
		Div(
			Img(
				Src(item.ImageURL),
				Alt(fmt.Sprintf("%s 이미지", item.KindLabel)),
				Width("70"),
				Height("70"),
				Loading("lazy"),
			),
		),
		Div(x.Show("$store.auth.isAuthed"), Class("closet-card__delete"),
			Attr("onclick", "event.stopPropagation();"),
			h.Delete(fmt.Sprintf("/items/%d", item.ID)),
			h.Target("closest article"),
			Attr("hx-confirm", "정말 삭제할까요?"),
			Attr("hx-on::after-request", "if(event.detail.successful){ this.closest('article').remove(); showInfo('아이템을 삭제했어요.'); }"),
			Text("X"),
		),
	)
}

// ItemDetailContent는 상세 다이얼로그 내용을 렌더링한다.
func ItemDetailContent(item ClosetItemDetail) Node {
	info := []Node{
		H5(Text(fmt.Sprintf("%s 상세 정보", item.KindLabel))),
		P(Class("caption"), Text(fmt.Sprintf("파일명: %s", item.Filename))),
	}
	if item.Dimension != "" {
		info = append(info, P(Class("caption"), Text(fmt.Sprintf("이미지 크기: %s", item.Dimension))))
	}
	if item.CreatedLabel != "" {
		info = append(info, P(Class("caption"), Text(fmt.Sprintf("업로드: %s", item.CreatedLabel))))
	}
	if item.TagLine != "" {
		info = append(info, P(Class("caption"), Text(item.TagLine)))
	}

	return Div(
		Div(Class("padding center-align"),
			Img(
				Src(item.ImageURL),
				Alt(fmt.Sprintf("%s 확대 이미지", item.KindLabel)),
				Width("220"),
				Loading("lazy"),
			),
			Group(info),
		),
		Div(Class("padding"),
			Form(
				Attr("hx-put", fmt.Sprintf("/items/%d", item.ID)),
				Attr("hx-indicator", "#item-detail-loading"),
				Attr("hx-on::after-request", "if(event.detail.successful){ closeModal('#item-detail-dialog'); showInfo('아이템 정보를 업데이트했어요. 임베딩을 다시 준비합니다.'); }"),
				h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '업데이트에 실패했어요. 다시 시도해 주세요.');"),
				Div(Class("field textarea label border"),
					Textarea(Name("meta_summary"), Rows("1"), Text(item.MetaSummary)),
					Label(Text("요약")),
					Span(Class("helper"), Text("간단한 설명을 작성해 주세요.")),
				),
				Div(Class("field label border"),
					Input(Type("text"), Name("meta_season"), Value(item.MetaSeason)),
					Label(Text("계절")),
				),
				Div(Class("field label border"),
					Input(Type("text"), Name("meta_style"), Value(item.MetaStyle)),
					Label(Text("스타일")),
				),
				Div(Class("field label border"),
					Input(Type("text"), Name("meta_colors"), Value(item.MetaColorsValue)),
					Label(Text("색상")),
					Span(Class("helper"), Text("콤마(,)로 여러 색상을 입력할 수 있어요.")),
				),
				Div(Class("field textarea label border"),
					Textarea(Name("tags"), Rows("3"), Text(item.TagsValue)),
					Label(Text("태그")),
					Span(Class("helper"), Text("콤마(,)나 줄바꿈으로 여러 태그를 입력할 수 있어요.")),
				),
				Div(Class("row"),
					Button(Class("button"), Type("submit"), Text("저장")),
					Button(Class("button border"), Type("button"), Attr("onclick", "closeModal('#item-detail-dialog')"), Text("닫기")),
				),
			),
		),
	)
}

func extractTags(raw interface{}) []string {
	var text string
	switch v := raw.(type) {
	case string:
		text = v
	case []byte:
		text = string(v)
	default:
		return nil
	}
	if strings.TrimSpace(text) == "" {
		return nil
	}
	parts := strings.Split(text, ",")
	tags := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}

func formatTagLine(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	sorted := make([]string, len(tags))
	copy(sorted, tags)
	sort.Strings(sorted)
	return "#" + strings.Join(sorted, " #")
}

func formatDimension(width sql.NullInt64, height sql.NullInt64) string {
	if !width.Valid || !height.Valid {
		return ""
	}
	return fmt.Sprintf("%dx%d", width.Int64, height.Int64)
}

func nullStringValue(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return strings.TrimSpace(v.String)
}
