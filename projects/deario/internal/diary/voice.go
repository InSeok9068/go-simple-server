package diary

import (
	"io"
	"log/slog"
	"net/http"

	aiclient "simple-server/internal/ai"
	"simple-server/pkg/util/authutil"

	"github.com/labstack/echo/v4"
)

// TranscribeDiaryVoice는 업로드된 음성을 텍스트로 변환한다.
func TranscribeDiaryVoice(c echo.Context) error {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return err
	}
	file, err := c.FormFile("audio")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "오디오 파일이 필요합니다.")
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	slog.Debug("음성 일기 변환", "user", uid, "size", len(data))
	text, err := aiclient.TranscribeAudio(c.Request().Context(), data, file.Header.Get("Content-Type"))
	if err != nil {
		slog.Error("음성 인식 실패", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "음성 인식에 실패했습니다.")
	}
	return c.String(http.StatusOK, text)
}
