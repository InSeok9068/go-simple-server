package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode/utf8"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views/components"

	"github.com/labstack/echo/v4"
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

	queries, err := db.GetQueries()
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

	items := make([]components.SearchResultItem, 0, len(diarys))
	for _, d := range diarys {
		items = append(items, components.SearchResultItem{
			Date:    d.Date,
			Snippet: snippetNodes(d.Content, q),
		})
	}

	return components.SearchResults(items).Render(c.Request().Context(), c.Response().Writer)
}

func snippetNodes(content, keyword string) components.SearchResultSnippet {
	lowerContent := strings.ToLower(content)
	lowerKeyword := strings.ToLower(keyword)
	byteIndex := strings.Index(lowerContent, lowerKeyword)

	contentRunes := []rune(content)
	contentRuneLen := len(contentRunes)

	if byteIndex == -1 {
		if contentRuneLen > 20 {
			return components.SearchResultSnippet{Prefix: string(contentRunes[:20]) + "..."}
		}
		return components.SearchResultSnippet{Prefix: content}
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

	return components.SearchResultSnippet{
		Prefix:   prefix,
		Match:    match,
		Suffix:   suffix,
		HasMatch: true,
	}
}
