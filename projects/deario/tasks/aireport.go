package tasks

import (
	"context"
	"fmt"
	"log/slog"
	aiclient "simple-server/internal/ai"
	"simple-server/internal/connection"
	"simple-server/projects/deario/db"
	"strings"
	"time"

	"maragu.dev/goqite"
	"maragu.dev/goqite/jobs"
)

var AiReportQ *goqite.Queue

// buildReportPrompt는 AI 상담 리포트 생성을 위한 프롬프트를 생성합니다.
func buildReportPrompt(diaries []db.Diary) string {
	var diaryEntries strings.Builder
	for _, diary := range diaries {
		// diary.Content가 비어있지 않은 경우에만 추가합니다.
		if strings.TrimSpace(diary.Content) != "" {
			diaryEntries.WriteString(fmt.Sprintf("날짜: %s\n내용: %s\n---\n", diary.Date, diary.Content))
		}
	}

	// 일기 내용이 없는 경우, 분석할 수 없다는 메시지를 포함한 프롬프트를 반환할 수 있습니다.
	if diaryEntries.Len() == 0 {
		return ""
	}

	return fmt.Sprintf(`당신은 따뜻하고 통찰력 있는 심리 상담사 '디어'입니다. 사용자의 지난 일기들을 바탕으로, 친구에게 말하듯 친근하고 다정한 어투로 심리 분석 리포트를 작성해주세요.

**당신의 역할:**
- 사용자의 감정, 생각, 관계의 패턴을 분석합니다.
- 비판이나 지적 대신, 공감과 지지를 표현합니다.
- 사용자가 스스로를 더 깊이 이해하고 긍정적인 방향으로 나아갈 수 있도록 돕습니다.
- 모든 결과는 마크다운 형식으로 작성해주세요.

**분석할 일기 내용:**
---
%s
---

**리포트 작성 구조:**

### 💌 디어의 마음 분석 리포트

OO님, 지난 한 달 동안의 소중한 마음 기록들을 제가 조심스럽게 읽어보았어요. OO님의 일상 속에서 느꼈던 다채로운 감정들을 함께 따라가 볼 수 있어 정말 의미 있는 시간이었어요. 아래에 제가 느낀 점들을 정리해 보았어요.

#### 1. 요즘 OO님의 마음 날씨는 어떤가요? 🌦️
(일기 전반에서 드러나는 핵심 감정 1~2개를 짚어주세요. 예를 들어, '설렘과 불안이 함께한 한 달이었네요' 또는 '차분하게 자신에게 집중하는 시간이 많았던 것 같아요' 와 같이 부드럽게 표현해주세요. 감정의 변화 추이가 있었다면 함께 언급해주세요.)

#### 2. OO님의 마음을 채우고 있는 생각들은 무엇인가요? 🤔
(가장 자주 언급되는 주제나 고민, 관심사를 요약해주세요. '새로운 프로젝트에 대한 기대감'이나 '친구 OOO와의 관계에 대한 고민'처럼 구체적인 내용을 바탕으로 작성해주세요. 반복되는 생각의 패턴이 있다면 짚어주세요.)

#### 3. 소중한 인연들과는 어떤 시간을 보냈나요? 👨‍👩‍👧‍👦
(일기에 등장하는 사람들과의 관계를 간략하게 분석해주세요. 긍정적/부정적 상호작용, 혹은 관계에 대한 고민 등을 따뜻한 시선으로 요약해주세요. 만약 관계 언급이 없다면, '이번 달은 온전히 자신에게 집중하는 시간을 보내신 것 같아요' 라고 언급해주세요.)

#### 💖 디어의 따뜻한 응원 한마디
(위 분석을 바탕으로, 사용자가 스스로를 칭찬해 줄 만한 부분이나 긍정적인 면을 찾아 강조해주세요. 그리고 앞으로 나아가는 데 도움이 될 만한 실천적인 조언이나 따뜻한 응원의 메시지를 한두 문장으로 전달해주세요. 예를 들어, '새로운 도전을 망설이면서도 한 발짝 내디딘 OO님의 용기가 정말 멋져요. 결과에 상관없이 그 과정 자체로도 충분히 의미 있답니다.' 와 같이 작성해주세요.)`, diaryEntries.String())
}

func GenerateAIReportJob() {
	apiReportDb, err := connection.AppDBOpen(false)
	if err != nil {
		slog.Error("데이터베이스 연결 실패", "error", err)
		return
	}
	AiReportQ = goqite.New(goqite.NewOpts{
		DB:   apiReportDb,
		Name: "ai-report",
	})
	r := jobs.NewRunner(jobs.NewRunnerOpts{
		Limit:        1,
		Log:          slog.Default(),
		PollInterval: 1 * time.Second,
		Extend:       5 * time.Minute, // 넉넉하게 5분의 타임아웃 설정
		Queue:        AiReportQ,
	})

	r.Register("ai-report", func(ctx context.Context, m []byte) error {
		uid := string(m)
		slog.Info("AI 리포트 생성 작업을 시작합니다.", "uid", uid)

		queries, err := db.GetQueries(ctx)
		if err != nil {
			slog.Error("쿼리 객체 생성 실패", "uid", uid, "error", err)
			return err
		}

		var allDiaries []db.Diary
		// 한 페이지에 7개의 일기를 가져오므로, 5페이지를 조회하여 최대 35개의 최근 일기를 가져옴
		for page := 1; page <= 5; page++ {
			// sqlc의 ListDiarys는 page 파라미터를 interface{}로 받으므로 page를 직접 전달
			diaries, err := queries.ListDiarys(ctx, db.ListDiarysParams{Uid: uid, Column2: int64(page)})
			if err != nil {
				slog.Error("일기 목록 조회 실패", "uid", uid, "page", page, "error", err)
				// 한 페이지 실패 시 다음 페이지 시도
				continue
			}
			if len(diaries) == 0 {
				// 더 이상 가져올 일기가 없으면 중단
				break
			}
			allDiaries = append(allDiaries, diaries...)
		}

		if len(allDiaries) == 0 {
			slog.Info("분석할 일기가 없어 작업을 종료합니다.", "uid", uid)
			return nil // 작업 성공 처리 (오류 아님)
		}

		slog.Info("일기 조회를 완료했으며, AI 프롬프트를 생성합니다.", "uid", uid, "diary_count", len(allDiaries))
		prompt := buildReportPrompt(allDiaries)
		if prompt == "" {
			slog.Info("내용이 있는 일기가 없어 작업을 종료합니다.", "uid", uid)
			return nil
		}

		// AI 클라이언트를 사용하여 리포트 생성 요청
		// 모델로 'gemini-2.5-pro'를 명시적으로 지정
		report, err := aiclient.Request(ctx, prompt, "gemini-2.5-pro")
		if err != nil {
			slog.Error("AI 리포트 생성 실패", "uid", uid, "error", err)
			return err // AI 요청 실패 시 재시도
		}

		// 생성된 리포트는 일단 로그로 출력합니다.
		// TODO: 향후 이 부분에 이메일 발송 또는 DB 저장 로직을 추가할 수 있습니다.
		slog.Info("AI 리포트 생성 성공", "uid", uid, "report", report)

		return nil
	})

	r.Start(context.Background())
}
