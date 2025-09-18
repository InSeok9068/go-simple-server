package handlers

import (
	"errors"
	"net/http"

	"simple-server/pkg/util/authutil"
	"simple-server/projects/portfolio/services"
	"simple-server/projects/portfolio/views"

	"github.com/labstack/echo/v4"
)

func CreateCategory(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateCategoryInput{
		Name:         c.FormValue("name"),
		Role:         c.FormValue("role"),
		ParentID:     c.FormValue("parent_id"),
		DisplayOrder: c.FormValue("display_order"),
	}

	if err := services.CreateCategory(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "카테고리를 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateSecurity(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateSecurityInput{
		Symbol:     c.FormValue("symbol"),
		Name:       c.FormValue("name"),
		Type:       c.FormValue("type"),
		CategoryID: c.FormValue("category_id"),
		Currency:   c.FormValue("currency"),
		Note:       c.FormValue("note"),
	}

	if err := services.CreateSecurity(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "종목을 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateAccount(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateAccountInput{
		CategoryID:     c.FormValue("category_id"),
		Name:           c.FormValue("name"),
		Provider:       c.FormValue("provider"),
		Balance:        c.FormValue("balance"),
		MonthlyContrib: c.FormValue("monthly_contrib"),
		Note:           c.FormValue("note"),
	}

	if err := services.CreateAccount(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "계좌를 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateHolding(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateHoldingInput{
		SecurityID:   c.FormValue("security_id"),
		Amount:       c.FormValue("amount"),
		TargetAmount: c.FormValue("target_amount"),
		Note:         c.FormValue("note"),
	}

	if err := services.CreateHolding(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "보유 종목을 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateAllocationTarget(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateAllocationTargetInput{
		CategoryID:   c.FormValue("category_id"),
		TargetWeight: c.FormValue("target_weight"),
		TargetAmount: c.FormValue("target_amount"),
		Note:         c.FormValue("note"),
	}

	if err := services.CreateAllocationTarget(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "리밸런싱 목표를 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateBudgetEntry(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateBudgetEntryInput{
		Name:      c.FormValue("name"),
		Direction: c.FormValue("direction"),
		Class:     c.FormValue("class"),
		Planned:   c.FormValue("planned"),
		Actual:    c.FormValue("actual"),
		Note:      c.FormValue("note"),
	}

	if err := services.CreateBudgetEntry(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "예산 항목을 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func CreateContributionPlan(c echo.Context) error {
	uid, err := requireUID(c)
	if err != nil {
		return err
	}

	input := services.CreateContributionPlanInput{
		SecurityID: c.FormValue("security_id"),
		Weight:     c.FormValue("weight"),
		Amount:     c.FormValue("amount"),
		Note:       c.FormValue("note"),
	}

	if err := services.CreateContributionPlan(c.Request().Context(), uid, input); err != nil {
		return handleServiceError(err, "적립 계획을 저장하지 못했습니다.")
	}

	return renderDashboard(c, uid)
}

func requireUID(c echo.Context) (string, error) {
	uid, err := authutil.SessionUID(c)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "로그인이 필요합니다.")
	}
	return uid, nil
}

func renderDashboard(c echo.Context, uid string) error {
	snapshot, err := services.LoadPortfolio(c.Request().Context(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "포트폴리오 데이터를 새로고침하지 못했습니다.")
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	return views.DashboardContent(snapshot).Render(c.Response().Writer)
}

func handleServiceError(err error, fallback string) error {
	var vErr *services.ValidationError
	if errors.As(err, &vErr) {
		return echo.NewHTTPError(http.StatusBadRequest, vErr.Message)
	}
	return echo.NewHTTPError(http.StatusInternalServerError, fallback)
}
