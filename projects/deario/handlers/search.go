package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/dateutil"
	"simple-server/projects/deario/db"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// SearchDiaries는 내용에서 키워드를 검색해 일기 목록을 반환한다.
func SearchDiaries(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	q := strings.TrimSpace(c.QueryParam("q"))
	if q == "" {
		return c.String(http.StatusOK, "")
	}

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diarys, err := queries.SearchDiarys(c.Request().Context(), db.SearchDiarysParams{
		Uid:     uid,
		Column2: sql.NullString{String: q, Valid: true},
	})
	if err != nil {
		return err
	}

	var lis []Node
	for _, d := range diarys {
		lis = append(lis,
			Li(
				A(Href(fmt.Sprintf("/?date=%s", d.Date)),
					Div(Text(dateutil.MustFormatDateKorSimpleWithWeekDay(d.Date))),
					Div(snippetNodes(d.Content, q)...),
				),
			),
		)
	}

	return Ul(Class("list"), Group(lis)).Render(c.Response().Writer)
}

func snippetNodes(content, keyword string) []Node {
	lower := strings.ToLower(content)
	lowerKey := strings.ToLower(keyword)
	idx := strings.Index(lower, lowerKey)
	if idx == -1 {
		if len(content) > 20 {
			content = content[:20] + "..."
		}
		return []Node{Text(content)}
	}

	start := idx - 5
	if start < 0 {
		start = 0
	}
	end := idx + len(keyword) + 5
	if end > len(content) {
		end = len(content)
	}

	prefix := content[start:idx]
	match := content[idx : idx+len(keyword)]
	suffix := content[idx+len(keyword) : end]

	if start > 0 {
		prefix = "..." + prefix
	}
	if end < len(content) {
		suffix += "..."
	}

	return []Node{Text(prefix), B(Text(match)), Text(suffix)}
}
