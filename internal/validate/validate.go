package validate

import (
	"log/slog"
	"reflect"

	localesko "github.com/go-playground/locales/ko"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translationsko "github.com/go-playground/validator/v10/translations/ko"
)

// 전역 검증기 및 번역기
var (
	Validate *validator.Validate
	Trans    ut.Translator
)

// Init는 go-playground/validator와 한국어 번역을 초기화한다.
func Init() {
	if Validate != nil && Trans != nil {
		return
	}

	// 한국어 로케일 설정
	ko := localesko.New()
	uni := ut.New(ko, ko)
	trans, found := uni.GetTranslator("ko")
	if !found {
		slog.Warn("한국어 번역기 초기화 실패, 기본 메시지 사용")
	}

	v := validator.New(validator.WithRequiredStructEnabled())

	// 필드명은 json 태그 사용(없으면 Go 필드명)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := fld.Tag.Get("json"); name != "" && name != "-" {
			return name
		}
		return fld.Name
	})

	// 기본 한국어 번역 등록(가능한 경우)
	if err := translationsko.RegisterDefaultTranslations(v, trans); err != nil {
		slog.Warn("검증기 한국어 번역 등록 실패", "error", err)
	}

	Validate = v
	Trans = trans
}

// ErrorMessage는 validator 오류를 한국어 메시지로 합쳐 반환한다.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	// ValidationErrors인 경우, 각 필드의 번역 메시지를 결합
	if verrs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			if Trans != nil {
				msgs = append(msgs, fe.Translate(Trans))
			} else {
				msgs = append(msgs, fe.Error())
			}
		}
		// 한 줄로 합치되, 다수 필드면 문장 구분을 위해 "; " 사용
		out := ""
		for i, m := range msgs {
			if i > 0 {
				out += "; "
			}
			out += m
		}
		return out
	}
	// 알 수 없는 오류는 기본 문자열 반환
	return err.Error()
}
