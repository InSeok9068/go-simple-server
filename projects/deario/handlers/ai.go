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
                Create a single image containing a 4-panel comic strip that tells a complete story without using any text, words, or written language. The four panels should be arranged in a single image, clearly separated but visually connected.

                For the image:
                1. Create a single image divided into 4 equal rectangular panels (2x2 grid)
                2. Each panel should be a self-contained illustration that flows naturally to the next
                3. Use only visual storytelling through composition, colors, lighting, and expressions
                4. Absolutely no text, captions, speech bubbles, signs, or written words
                5. Show clear character emotions and actions to convey the narrative
                6. Maintain visual consistency across all panels
                7. Each panel should focus on a single, meaningful moment or emotion

                The comic should tell its story through pure visual language, like a wordless graphic novel. The sequence of four panels should show a clear beginning, development, and resolution of a simple, relatable situation or emotion.

                Use the context provided in the 'contents' field only as inspiration for the mood and setting, but do not include any text or literal elements from it in the image.`
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

               ※ 감정을 깊게 공감하고 나서 %s

               이해했다는말이나 이런거 하지말고 바로 답변해줘

               응답은 반드시 Markdown 형식으로 해줘.
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

	return nil
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

	return nil
}
