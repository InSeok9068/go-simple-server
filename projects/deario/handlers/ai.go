package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	aiclient "simple-server/internal/ai"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/tasks"

	"github.com/labstack/echo/v4"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"maragu.dev/goqite/jobs"
)

// GenerateAIFeedback는 일기 내용을 기반으로 AI 피드백이나 이미지를 생성한다.
func GenerateAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	content := c.FormValue("content")
	typeValue := c.QueryParam("type")

	if content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "내용을 입력해주세요.")
	}

	slog.Debug("AI 피드백", "user", uid, "content", content, "type", typeValue)

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
							Create a vertical (1x4) comic strip in a single image
							The image should contain 4 equal rectangular panels arranged vertically from top to bottom.

							Requirements
							- Divide a single image into 4 equal vertical rectangular panels, arranged from top to bottom (1x4 layout).
							- Use only visual storytelling — composition, colors, lighting, and facial expressions should convey the story and emotions.
							- No text, captions, speech bubbles, signs, or any written language at all.
							- Maintain visual consistency across all panels.
							- Each panel should focus on one meaningful moment or emotion.
							- Image size: height 700px, width 320px.
							- The four panels together must tell a complete story.
							- The story should have a clear beginning (introduction), development, conflict or change, and resolution.
							- The situation or emotion should be simple, relatable, and easy to understand without words.`
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
	}

	prompt := fmt.Sprintf(`아래의 내용은 나의 오늘 하루의 일기야
							내용 : %s

              ※ 유의사항
							1. 감정을 공감하고 나서 %s 이해했다는말이나 이런거 하지말고 바로 답변해줘.
							2. 답변 길이는 250자 ~ 500자 사이로 해줘.
							3. 답변은 너무 오바하지말고 담백하면서도 친절하게 해줘.
							4. 응답은 반드시 Markdown 형식으로 해줘.
              `, content, typeStr)
	result, err := aiclient.Request(c.Request().Context(), prompt)
	if err != nil {
		return err
	}

	return Div(
		Input(Type("hidden"), Name("ai-feedback"), Value(result)),
		Div(ID("ai-feedback-markdown"), Text(result)),
	).Render(c.Response().Writer)
}

// SaveAIFeedback는 생성된 AI 피드백과 이미지를 저장한다.
func SaveAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.FormValue("date")
	aiFeedback := c.FormValue("ai-feedback")
	aiImage := c.FormValue("ai-image")

	queries, err := db.GetQueries()
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

	if err := queries.UpdateDiaryOfAiFeedback(c.Request().Context(), db.UpdateDiaryOfAiFeedbackParams{
		ID:         diary.ID,
		AiFeedback: aiFeedback,
		AiImage:    aiImage,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기요정 저장에 실패하였습니다.")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAIFeedback는 저장된 피드백이나 이미지를 반환한다.
func GetAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date := c.QueryParam("date")

	queries, err := db.GetQueries()
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

	if diary.AiImage != "" {
		return Div(
			Input(Type("hidden"), Name("ai-image"), Value(diary.AiImage)),
			Img(Style("width:320px"), Src(fmt.Sprintf("data:image/png;base64,%s", diary.AiImage))),
		).Render(c.Response().Writer)
	}

	if diary.AiFeedback == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "저장된 일기요정이 없습니다.")
	}

	return Div(
		Input(Type("hidden"), Name("ai-feedback"), Value(diary.AiFeedback)),
		Div(ID("ai-feedback-markdown"), Text(diary.AiFeedback)),
	).Render(c.Response().Writer)
}

// GenerateAIReport AI 상담 리포트를 생성한다.
func GenerateAIReport(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	// 큐에 작업 추가
	if err := jobs.Create(c.Request().Context(), tasks.AiReportQ, "ai-report", []byte(uid)); err != nil {
		slog.Error("AI 리포트 발송 실패", "error", err)
	}

	return c.NoContent(http.StatusAccepted)
}
