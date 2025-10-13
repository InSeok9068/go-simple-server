package validate

import (
	"log/slog"
	"reflect"
	"strings"

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
	return ErrorMessageOf(nil, err)
}

// target(검증했던 원본 구조체 or 포인터)을 주면 그 안의 `message` 태그를 우선 적용
func ErrorMessageOf(target any, err error) string {
	if err == nil {
		return ""
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	msgs := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		// 1) message 태그 우선 (심플하게: 최상위 구조체의 바로 그 필드만)
		if m, ok := lookupFieldMessageTag(target, fe); ok && m != "" {
			msgs = append(msgs, m)
			continue
		}
		// 2) 번역(ko) 있으면 번역, 없으면 기본 에러
		if Trans != nil {
			msgs = append(msgs, fe.Translate(Trans))
		} else {
			msgs = append(msgs, fe.Error())
		}
	}
	return joinMessages(msgs)
}

// 최상위 구조체에서 해당 필드의 `message` 태그를 찾아 반환
func lookupFieldMessageTag(target any, fe validator.FieldError) (string, bool) {
	if target == nil {
		return "", false
	}
	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", false
	}

	// v10에선 Top/Parent가 없으므로, 간단히 최상위 구조체의 필드만 본다.
	// (중첩 구조는 필요해지면 확장)
	// fe.StructField()는 Go 필드명 반환 (JSON 이름 아님)
	if fld, ok := t.FieldByName(fe.StructField()); ok {
		if msg := fld.Tag.Get("message"); strings.TrimSpace(msg) != "" {
			return msg, true
		}
	}
	return "", false
}

func joinMessages(msgs []string) string {
	switch len(msgs) {
	case 0:
		return ""
	case 1:
		return msgs[0]
	default:
		out := msgs[0]
		for _, m := range msgs[1:] {
			out += "; " + m
		}
		return out
	}
}
