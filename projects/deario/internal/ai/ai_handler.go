package ai

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	aiclient "simple-server/internal/ai"
	"simple-server/pkg/util/authutil"
	"simple-server/projects/deario/db"
	"simple-server/projects/deario/internal/deariodate"
	"simple-server/projects/deario/internal/notification"
	"simple-server/projects/deario/views/components"

	"github.com/labstack/echo/v4"
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

	typeStr, err := aiFeedbackInstruction(typeValue)
	if err != nil {
		return err
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
		return components.AiImageResult(result).Render(c.Request().Context(), c.Response().Writer)
	}

	prompt := fmt.Sprintf(`너는 사용자의 일기를 읽고 짧고 다정한 코멘트를 남겨주는 "일기요정"이야.

아래 일기를 읽고, 일기 속 상황이나 표현 하나를 가볍게 짚은 뒤 사용자가 느꼈을 마음을 담백하게 말해줘.
너무 감성적인 말, 과한 위로, 특별한 관계인 척하는 표현은 피하고 자연스럽게 답해줘.

오늘의 일기:
%s

응답 방향:
- 요청한 방식: %s
- 일기 속 구체적인 장면이나 표현을 하나만 짚어줘.
- 감정은 단정하지 말고 "~했을 수도 있겠다"처럼 조심스럽게 말해줘.
- "나는 네 곁에 있어", "언제나", "소중한 마음"처럼 과하게 다정한 표현은 피해줘.
- 필요하면 오늘 바로 해볼 수 있는 작은 제안 하나만 덧붙여줘.
- 마지막 문장은 짧은 여운이나 질문 하나로 끝내줘.
- Markdown 형식으로, 제목 없이 2~3개의 짧은 문단으로 작성해줘.
- 답변은 300자에서 500자 사이로 해줘.
`, content, typeStr)
	// 상황별 감정 해석을 더 강화하고 싶을 때 교체할 프롬프트:
	//
	// 너는 사용자의 일기를 읽고 다정하게 답장해주는 "일기요정"이야.
	//
	// 아래 일기를 읽고, 오늘 사용자가 어떤 감정을 느꼈는지 상황별로 부드럽게 짚어줘.
	// 감정을 단정하지 말고 "~했을 수도 있겠다", "~하게 느껴졌을 것 같아"처럼 조심스럽게 말해줘.
	// 과장된 표현, 훈계, 진단처럼 들리는 말은 피하고, 친구처럼 담백하고 따뜻하게 답해줘.
	//
	// 오늘의 일기:
	// %s
	//
	// 응답 방향:
	// - 요청한 방식: %s
	// - 첫 문장은 "이해했어", "알겠어" 같은 말로 시작하지 말고 바로 공감으로 시작해줘.
	// - 일기 안에서 의미 있는 상황을 1~3개 정도 자연스럽게 짚어줘.
	// - 각 상황에서 느꼈을 법한 감정과 그 감정이 생긴 이유를 짧게 연결해줘.
	// - 상황별 설명 뒤에는 필요한 경우 작은 위로나 현실적인 제안을 덧붙여줘.
	// - 충고가 필요한 경우에도 단정하지 말고, 사용자가 오늘 바로 해볼 수 있는 작은 제안으로 말해줘.
	// - 답변 길이는 350자에서 600자 사이로 해줘.
	// - Markdown 형식으로 작성하되, 제목은 쓰지 말고 짧은 문단과 필요한 경우 목록 1개만 사용해줘.
	result, err := aiclient.Request(c.Request().Context(), prompt)
	if err != nil {
		return err
	}

	return components.AiFeedbackResult(result).Render(c.Request().Context(), c.Response().Writer)
}

// SaveAIFeedback는 생성된 AI 피드백과 이미지를 저장한다.
func SaveAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date, err := deariodate.NormalizeRequired(c.FormValue("date"))
	if err != nil {
		return err
	}
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

// DeleteAIFeedback는 저장된 AI 피드백과 이미지를 삭제한다.
func DeleteAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date, err := deariodate.NormalizeRequired(c.FormValue("date"))
	if err != nil {
		return err
	}

	queries, err := db.GetQueries()
	if err != nil {
		return err
	}

	diary, err := queries.GetDiary(c.Request().Context(), db.GetDiaryParams{
		Uid:  uid,
		Date: date,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNoContent)
		}
		return err
	}

	diary.AiFeedback = ""
	diary.AiImage = ""
	if diary.Content == "" && !hasDiaryImageData(diary) {
		if err := queries.DeleteDiary(c.Request().Context(), diary.ID); err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	}

	if err := queries.UpdateDiaryOfAiFeedback(c.Request().Context(), db.UpdateDiaryOfAiFeedbackParams{
		ID:         diary.ID,
		AiFeedback: "",
		AiImage:    "",
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "일기요정 삭제에 실패하였습니다.")
	}

	return c.NoContent(http.StatusNoContent)
}

// GetAIFeedback는 저장된 피드백이나 이미지를 반환한다.
func GetAIFeedback(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	date, err := deariodate.NormalizeRequired(c.QueryParam("date"))
	if err != nil {
		return err
	}

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
		return components.AiImageResult(diary.AiImage).Render(c.Request().Context(), c.Response().Writer)
	}

	if diary.AiFeedback == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "저장된 일기요정이 없습니다.")
	}

	return components.AiFeedbackResult(diary.AiFeedback).Render(c.Request().Context(), c.Response().Writer)
}

func hasDiaryImageData(d db.Diary) bool {
	return d.ImageUrl1 != "" || d.ImageUrl2 != "" || d.ImageUrl3 != ""
}

func aiFeedbackInstruction(typeValue string) (string, error) {
	switch typeValue {
	case "1":
		return "오늘 일기에서 사용자가 애쓴 부분을 찾아 따뜻하게 칭찬해줘", nil
	case "2":
		return "사용자가 느꼈을 마음을 충분히 받아주고 부드럽게 위로해줘", nil
	case "3":
		return "단정하거나 훈계하지 말고, 도움이 될 만한 작은 조언을 하나만 건네줘", nil
	case "4":
		return `
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
							- The situation or emotion should be simple, relatable, and easy to understand without words.`, nil
	default:
		return "", echo.NewHTTPError(http.StatusBadRequest, "일기요정 요청 유형이 올바르지 않습니다.")
	}
}

// GenerateAIReport AI 상담 리포트를 생성한다.
func GenerateAIReport(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}

	// 큐에 작업 추가
	if err := notification.EnqueueAIReport(c.Request().Context(), uid); err != nil {
		slog.Error("AI 리포트 발송 실패", "error", err)
	}

	return c.NoContent(http.StatusAccepted)
}
