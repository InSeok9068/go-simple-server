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
	return Div(ID("items-list"), Class("stack gap-lg"), Group(sections))
}

func itemGroup(kind string, items []ClosetItem) Node {
	title := KindLabel(kind)
	total := len(items)

	cardNodes := make([]Node, 0, total)
	for _, item := range items {
		cardNodes = append(cardNodes, itemCard(item))
	}

	return Section(Class("stack gap-sm"),
		Nav(Class("between align-center"),
			H6(Class("title"), Text(title)),
			Span(Class("caption muted"), Text(fmt.Sprintf("%d개", total))),
		),
		If(total == 0,
			Div(Class("surface-container"),
				P(Class("caption"), Text("조건에 맞는 아이템이 아직 없어요.")),
			),
		),
		If(total > 0,
			Div(Class("row scroll gap-sm"), Group(cardNodes)),
		),
	)
}

func itemCard(item ClosetItem) Node {
	return Article(Class("card border relative"),
		Div(Class("center"),
			Img(
				Src(item.ImageURL),
				Alt(fmt.Sprintf("%s 이미지", item.KindLabel)),
				Width("150"),
				Height("150"),
				Loading("lazy"),
			),
		),
		Div(x.Show("$store.auth.isAuthed"), Class("closet-card__delete"),
			h.Delete(fmt.Sprintf("/items/%d", item.ID)),
			h.Target("closest article"),
			Attr("hx-confirm", "정말 삭제할까요?"),
			Attr("hx-on::after-request", "if(event.detail.successful){ this.closest('article').remove(); showInfo('아이템을 삭제했어요.'); }"),
			I(Class("icon"), Text("close")),
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
