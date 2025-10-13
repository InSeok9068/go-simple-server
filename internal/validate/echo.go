package validate

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

// EchoValidator는 Echo의 Validator 인터페이스 구현체다.
type EchoValidator struct{}

func NewEchoValidator() *EchoValidator {
    Init()
    return &EchoValidator{}
}

// Validate는 구조체 태그 기반 검증을 수행한다.
func (ev *EchoValidator) Validate(i interface{}) error {
    Init()
    return Validate.Struct(i)
}

// HTTPError는 검증 오류를 422 상태코드의 echo.HTTPError로 변환한다.
func HTTPError(err error) *echo.HTTPError {
    return echo.NewHTTPError(http.StatusUnprocessableEntity, ErrorMessage(err))
}

