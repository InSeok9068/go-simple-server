package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	aiclient "simple-server/internal/ai"
	"simple-server/pkg/util/authutil"
	"simple-server/pkg/util/dateutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/views"
	shared "simple-server/shared/views"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		date = strings.ReplaceAll(date, "-", "")
	}
	return views.Index(os.Getenv("APP_TITLE"), date).Render(c.Response().Writer)
}

func Login(c echo.Context) error {
	return shared.Login().Render(c.Response().Writer)
}

func Diary(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		date = strings.ReplaceAll(date, "-", "")
	}

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		return views.DiaryContentForm(date, "").Render(c.Response().Writer)
	} else {
		return views.DiaryContentForm(diary.Date, diary.Content).Render(c.Response().Writer)
	}
}

func DiaryList(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diarys, _ := queries.ListDiarys(c.Request().Context(), db.ListDiarysParams{
		Uid:     uid,
		Column2: page,
	})

	var lis []Node
	for _, diary := range diarys {
		lis = append(lis,
			Li(
				A(Href(fmt.Sprintf("/?date=%s", diary.Date)),
					Text(dateutil.MustFormatDateKorWithWeekDay(diary.Date)),
				),
			),
		)
	}

	return Group(lis).Render(c.Response().Writer)
}

func DiaryRandom(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diary, err := queries.GetDiaryRandom(c.Request().Context(), uid)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "작성한 일기장이 없습니다.")
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`<script>location.href = "/?date=%s";</script>`, diary.Date))
}

func Save(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	content := c.FormValue("content")

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		_, err = queries.CreateDiary(c.Request().Context(), db.CreateDiaryParams{
			Uid:     uid,
			Content: content,
			Date:    date,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "등록 실패")
		}
	} else {
		if content == "" {
			err = queries.DeleteDiary(c.Request().Context(), diary.ID)
		} else {
			_, err = queries.UpdateDiary(c.Request().Context(), db.UpdateDiaryParams{
				Content: content,
				ID:      diary.ID,
			})
		}

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "수정 실패")
		}
	}

	return nil
}

func AiFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	content := c.FormValue("content")
	typeValue := c.QueryParam("type")

	slog.Info("AiFeedback", "user", uid, "content", content, "type", typeValue)

	var typeStr string
	switch typeValue {
	case "1":
		typeStr = "칭찬을 해줘"
	case "2":
		typeStr = "위로를 해줘"
	case "3":
		typeStr = "충고를 해줘"
	case "4":
		typeStr = `
		Create a single illustrated image that represents an emotional or meaningful moment from someone's day.
		
		Do not include any text, captions, signs, labels, handwriting, or words in any language in the image.  
		Absolutely no text or characters should appear anywhere.
		
		Just imagine a beautiful or reflective scene like a page from a wordless picture book, using only colors, lighting, composition, and elements to tell the story.
		
		You will be given some context in the 'contents' field, but this is for inspiration only.  
		**Do not copy, draw, or include any part of the content text in the image.**  
		Use it only to understand the mood and setting.
		
		Only return an image. Do not generate or show any text in the picture.`
	}

	if typeValue == "4" {
		prompt := fmt.Sprintf(`
		%s

		content : %s
		`, typeStr, content)
		result, err := aiclient.ImageRequest(c.Request().Context(), prompt)
		if err != nil {
			return err
		}
		return Div(
			Input(Type("hidden"), Name("ai-image"), Value(result)),
			Img(Style("width:320px"), Src(fmt.Sprintf("data:image/png;base64,%s", result))),
		).Render(c.Response().Writer)
	} else {
		prompt := fmt.Sprintf(`아래의 내용은 나의 오늘 하루의 일기야
		내용 : %s
	
		※ 감정을 깊게 공감하고 나서 %s
		
		이해했다는말이나 이런거 하지말고 바로 답변해줘
		
		[응답형태는 마크다운이 아닌 <textarea>에 붙여넣을거라서 텍스트에 띄어쓰기나 줄바꿈으로 가독성을 높여줘]
		`, content, typeStr)
		result, err := aiclient.Request(c.Request().Context(), prompt)
		if err != nil {
			return err
		}

		return Div(
			Input(Type("hidden"), Name("ai-feedback"), Value(result)),
			Text(result),
		).Render(c.Response().Writer)
	}
}

func AiFeedbackSave(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	aiFeedback := c.FormValue("ai-feedback")
	aiImage := c.FormValue("ai-image")

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "작성한 일기가 없습니다.")
	}

	err = queries.UpdateDiaryOfAiFeedback(c.Request().Context(), db.UpdateDiaryOfAiFeedbackParams{
		ID:         diary.ID,
		Aifeedback: aiFeedback,
		Aiimage:    aiImage,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기요정 저장에 실패하였습니다.")
	}

	return nil
}

func GetAiFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.QueryParam("date")

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "저장된 일기가 없습니다.")
	}

	if diary.Aiimage != "" {
		return Div(
			Input(Type("hidden"), Name("ai-image"), Value(diary.Aiimage)),
			Img(Style("width:320px"), Src(fmt.Sprintf("data:image/png;base64,%s", diary.Aiimage))),
		).Render(c.Response().Writer)
	}

	if diary.Aifeedback == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "저장된 일기요정이 없습니다.")
	}

	return Div(
		Input(Type("hidden"), Name("ai-feedback"), Value(diary.Aifeedback)),
		Text(diary.Aifeedback),
	).Render(c.Response().Writer)
}

func SavePushKey(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return err
	}

	token := data["token"].(string)

	queries, err := db.GetQueries(c.Request().Context())
	if err != nil {
		return err
	}

	_, err = queries.GetPushKey(c.Request().Context(), uid)
	if err != nil {
		_ = queries.CreatePushKey(c.Request().Context(), db.CreatePushKeyParams{
			Uid:   uid,
			Token: token,
		})
	} else {
		_ = queries.UpdatePushKey(c.Request().Context(), db.UpdatePushKeyParams{
			Uid:   uid,
			Token: token,
		})
	}

	return nil
}
