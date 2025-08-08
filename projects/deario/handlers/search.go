package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

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
	lowerContent := strings.ToLower(content)
	lowerKeyword := strings.ToLower(keyword)
	byteIndex := strings.Index(lowerContent, lowerKeyword)

	contentRunes := []rune(content)
	contentRuneLen := len(contentRunes)

	if byteIndex == -1 {
		if contentRuneLen > 20 {
			return []Node{Text(string(contentRunes[:20]) + "...")}
		}
		return []Node{Text(content)}
	}

	runeIndex := utf8.RuneCountInString(lowerContent[:byteIndex])
	keywordRuneLen := len([]rune(keyword))

	startRune := runeIndex - 5
	if startRune < 0 {
		startRune = 0
	}
	endRune := runeIndex + keywordRuneLen + 5
	if endRune > contentRuneLen {
		endRune = contentRuneLen
	}

	prefix := string(contentRunes[startRune:runeIndex])
	match := string(contentRunes[runeIndex : runeIndex+keywordRuneLen])
	suffix := string(contentRunes[runeIndex+keywordRuneLen : endRune])

	if startRune > 0 {
		prefix = "..." + prefix
	}
	if endRune < contentRuneLen {
		suffix += "..."
	}

	return []Node{Text(prefix), B(Text(match)), Text(suffix)}
}
