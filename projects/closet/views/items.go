package views

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"simple-server/projects/closet/db"
)

var kindLabels = map[string]string{
	"top":       "상의",
	"bottom":    "하의",
	"shoes":     "신발",
	"accessory": "액세서리",
}

var kindRenderOrder = []string{"top", "bottom", "shoes", "accessory"}

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
