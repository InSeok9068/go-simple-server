package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"simple-server/projects/portfolio/db"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type CreateCategoryInput struct {
	Name         string
	Role         string
	ParentID     string
	DisplayOrder string
}

type CreateSecurityInput struct {
	Symbol     string
	Name       string
	Type       string
	CategoryID string
	Currency   string
	Note       string
}

type CreateAccountInput struct {
	CategoryID     string
	Name           string
	Provider       string
	Balance        string
	MonthlyContrib string
	Note           string
}

type CreateHoldingInput struct {
	SecurityID   string
	Amount       string
	TargetAmount string
	Note         string
}

type CreateAllocationTargetInput struct {
	CategoryID   string
	TargetWeight string
	TargetAmount string
	Note         string
}

type CreateBudgetEntryInput struct {
	Name      string
	Direction string
	Class     string
	Planned   string
	Actual    string
	Note      string
}

type CreateContributionPlanInput struct {
	SecurityID string
	Weight     string
	Amount     string
	Note       string
}

func CreateCategory(ctx context.Context, uid string, input CreateCategoryInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &ValidationError{Message: "카테고리명을 입력해 주세요."}
	}

	role := strings.TrimSpace(input.Role)
	if !isValidCategoryRole(role) {
		return &ValidationError{Message: "카테고리 유형을 선택해 주세요."}
	}

	displayOrder, err := parseIntField(input.DisplayOrder, "정렬 순서는 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	parentID := strings.TrimSpace(input.ParentID)
	var parent sql.NullString
	if parentID != "" {
		parent = sql.NullString{String: parentID, Valid: true}
	}

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateCategory(ctx, db.CreateCategoryParams{
			Uid:          uid,
			Name:         name,
			Role:         role,
			ParentID:     parent,
			DisplayOrder: displayOrder,
		}); err != nil {
			return fmt.Errorf("카테고리 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateSecurity(ctx context.Context, uid string, input CreateSecurityInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &ValidationError{Message: "종목명을 입력해 주세요."}
	}

	stype := strings.TrimSpace(input.Type)
	if !isValidSecurityType(stype) {
		return &ValidationError{Message: "종목 유형을 선택해 주세요."}
	}

	symbol := strings.TrimSpace(input.Symbol)
	currency := strings.TrimSpace(input.Currency)
	if currency == "" {
		currency = "KRW"
	}

	categoryID := strings.TrimSpace(input.CategoryID)
	var category sql.NullString
	if categoryID != "" {
		category = sql.NullString{String: categoryID, Valid: true}
	}

	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateSecurity(ctx, db.CreateSecurityParams{
			Uid:        uid,
			Symbol:     symbol,
			Name:       name,
			Type:       stype,
			CategoryID: category,
			Currency:   currency,
			Note:       note,
		}); err != nil {
			return fmt.Errorf("종목 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateAccount(ctx context.Context, uid string, input CreateAccountInput) error {
	categoryID := strings.TrimSpace(input.CategoryID)
	if categoryID == "" {
		return &ValidationError{Message: "계좌 카테고리를 선택해 주세요."}
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &ValidationError{Message: "계좌명을 입력해 주세요."}
	}

	balance, err := parseIntField(input.Balance, "잔액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	monthly, err := parseIntField(input.MonthlyContrib, "월 납입액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	provider := strings.TrimSpace(input.Provider)
	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateAccount(ctx, db.CreateAccountParams{
			Uid:            uid,
			CategoryID:     categoryID,
			Name:           name,
			Provider:       provider,
			Balance:        balance,
			MonthlyContrib: monthly,
			Note:           note,
		}); err != nil {
			return fmt.Errorf("계좌 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateHolding(ctx context.Context, uid string, input CreateHoldingInput) error {
	securityID := strings.TrimSpace(input.SecurityID)
	if securityID == "" {
		return &ValidationError{Message: "보유 종목을 선택해 주세요."}
	}

	amount, err := parseIntField(input.Amount, "현재 금액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	target, err := parseIntField(input.TargetAmount, "목표 금액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateHolding(ctx, db.CreateHoldingParams{
			Uid:          uid,
			SecurityID:   securityID,
			Amount:       amount,
			TargetAmount: target,
			Note:         note,
		}); err != nil {
			return fmt.Errorf("보유 종목 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateAllocationTarget(ctx context.Context, uid string, input CreateAllocationTargetInput) error {
	categoryID := strings.TrimSpace(input.CategoryID)
	if categoryID == "" {
		return &ValidationError{Message: "목표 카테고리를 선택해 주세요."}
	}

	weight, amount, err := parseWeightAndAmount(
		input.TargetWeight,
		"목표 비중은 숫자로 입력해 주세요.",
		"목표 비중은 0~100 사이로 입력해 주세요.",
		input.TargetAmount,
		"목표 금액은 숫자로 입력해 주세요.",
	)
	if err != nil {
		return err
	}

	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateAllocationTarget(ctx, db.CreateAllocationTargetParams{
			Uid:          uid,
			CategoryID:   categoryID,
			TargetWeight: weight,
			TargetAmount: amount,
			Note:         note,
		}); err != nil {
			return fmt.Errorf("리밸런싱 목표 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateBudgetEntry(ctx context.Context, uid string, input CreateBudgetEntryInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &ValidationError{Message: "예산 항목명을 입력해 주세요."}
	}

	direction := strings.TrimSpace(input.Direction)
	if direction != "income" && direction != "expense" {
		return &ValidationError{Message: "수입/지출 구분을 선택해 주세요."}
	}

	class := strings.TrimSpace(input.Class)
	if class != "fixed" && class != "variable" {
		return &ValidationError{Message: "고정/유동 구분을 선택해 주세요."}
	}

	planned, err := parseIntField(input.Planned, "계획 금액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	actual, err := parseIntField(input.Actual, "실제 금액은 숫자로 입력해 주세요.")
	if err != nil {
		return err
	}

	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateBudgetEntry(ctx, db.CreateBudgetEntryParams{
			Uid:       uid,
			Name:      name,
			Direction: direction,
			Class:     class,
			Planned:   planned,
			Actual:    actual,
			Note:      note,
		}); err != nil {
			return fmt.Errorf("예산 항목 저장 실패: %w", err)
		}
		return nil
	})
}

func CreateContributionPlan(ctx context.Context, uid string, input CreateContributionPlanInput) error {
	securityID := strings.TrimSpace(input.SecurityID)
	if securityID == "" {
		return &ValidationError{Message: "적립할 종목을 선택해 주세요."}
	}

	weight, amount, err := parseWeightAndAmount(
		input.Weight,
		"적립 비중은 숫자로 입력해 주세요.",
		"적립 비중은 0~100 사이로 입력해 주세요.",
		input.Amount,
		"적립 금액은 숫자로 입력해 주세요.",
	)
	if err != nil {
		return err
	}

	note := strings.TrimSpace(input.Note)

	return withQueries(func(queries *db.Queries) error {
		if err := queries.CreateContributionPlan(ctx, db.CreateContributionPlanParams{
			Uid:        uid,
			SecurityID: securityID,
			Weight:     weight,
			Amount:     amount,
			Note:       note,
		}); err != nil {
			return fmt.Errorf("적립 계획 저장 실패: %w", err)
		}
		return nil
	})
}

func parseIntField(raw string, message string) (int64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, &ValidationError{Message: message}
	}
	return value, nil
}

func parseFloatField(raw string, message string) (float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, &ValidationError{Message: message}
	}
	return value, nil
}

func parsePercentageField(raw string, parseMessage string, rangeMessage string) (float64, error) {
	value, err := parseFloatField(raw, parseMessage)
	if err != nil {
		return 0, err
	}
	if value < 0 || value > 100 {
		return 0, &ValidationError{Message: rangeMessage}
	}
	return value, nil
}

func parseWeightAndAmount(weightRaw string, weightParseMessage string, weightRangeMessage string, amountRaw string, amountMessage string) (float64, int64, error) {
	weight, err := parsePercentageField(weightRaw, weightParseMessage, weightRangeMessage)
	if err != nil {
		return 0, 0, err
	}

	amount, err := parseIntField(amountRaw, amountMessage)
	if err != nil {
		return 0, 0, err
	}

	return weight, amount, nil
}

func isValidCategoryRole(role string) bool {
	switch role {
	case "liquidity", "growth", "income", "protection", "other":
		return true
	default:
		return false
	}
}

func withQueries(action func(*db.Queries) error) error {
	queries, err := db.GetQueries()
	if err != nil {
		return fmt.Errorf("쿼리 초기화 실패: %w", err)
	}
	return action(queries)
}

func isValidSecurityType(t string) bool {
	switch t {
	case "stock", "etf", "fund", "bond", "crypto", "other":
		return true
	default:
		return false
	}
}
