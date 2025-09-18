package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"simple-server/projects/portfolio/db"
)

const DemoUID = "uid_demo_portfolio"

type PortfolioSnapshot struct {
	UID                 string
	TotalAsset          int64
	Categories          []CategorySummary
	Securities          []SecuritySummary
	CategoryAllocations []CategoryAllocation
	Accounts            []Account
	Holdings            []Holding
	RebalanceGaps       []RebalanceGap
	BudgetEntries       []BudgetEntry
	ContributionPlans   []ContributionPlan
	IsDemo              bool
}

type CategorySummary struct {
	ID           string
	Name         string
	Role         string
	ParentID     string
	DisplayOrder int64
}

type CategoryAllocation struct {
	CategoryID   string
	CategoryName string
	Amount       int64
	WeightPct    float64
}

type SecuritySummary struct {
	ID         string
	Name       string
	Symbol     string
	Type       string
	CategoryID string
	Currency   string
	Note       string
}

type Account struct {
	ID             string
	Name           string
	Provider       string
	CategoryName   string
	Balance        int64
	MonthlyContrib int64
	Note           string
}

type Holding struct {
	ID             string
	SecurityID     string
	SecurityName   string
	SecuritySymbol string
	SecurityType   string
	CategoryName   string
	Currency       string
	Amount         int64
	TargetAmount   int64
	Note           string
}

type RebalanceGap struct {
	CategoryID     string
	CategoryName   string
	TargetWeight   float64
	TargetAmount   int64
	CurrentAmount  int64
	TargetByWeight int64
	TargetFinal    int64
	GapAmount      int64
}

type BudgetEntry struct {
	ID        string
	Name      string
	Direction string
	Class     string
	Planned   int64
	Actual    int64
	Note      string
}

type ContributionPlan struct {
	ID             string
	SecurityID     string
	SecurityName   string
	SecuritySymbol string
	CategoryName   string
	Weight         float64
	Amount         int64
	Note           string
}

func LoadPortfolio(ctx context.Context, uid string) (*PortfolioSnapshot, error) {
        queries, err := db.GetQueries()
        if err != nil {
                return nil, fmt.Errorf("쿼리 초기화 실패: %w", err)
        }

        snapshot := &PortfolioSnapshot{UID: uid, IsDemo: uid == DemoUID}

        total, err := queries.GetTotalAsset(ctx, uid)
        if err != nil {
                if !errors.Is(err, sql.ErrNoRows) {
                        return nil, fmt.Errorf("총자산 조회 실패: %w", err)
                }
        } else {
                snapshot.TotalAsset = total.TotalAmount
        }

        fillers := []func(context.Context, *db.Queries, *PortfolioSnapshot) error{
                fillCategories,
                fillSecurities,
                fillCategoryAllocations,
                fillAccounts,
                fillHoldings,
                fillRebalanceGaps,
                fillBudgetEntries,
                fillContributionPlans,
        }
        for _, filler := range fillers {
                if err := filler(ctx, queries, snapshot); err != nil {
                        return nil, err
                }
        }

        return snapshot, nil
}

func fillCategories(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListCategories(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("카테고리 목록 조회 실패: %w", err)
	}

	snapshot.Categories = make([]CategorySummary, 0, len(rows))
	for _, row := range rows {
		snapshot.Categories = append(snapshot.Categories, CategorySummary{
			ID:           row.ID,
			Name:         row.Name,
			Role:         row.Role,
			ParentID:     nullString(row.ParentID),
			DisplayOrder: row.DisplayOrder,
		})
	}
	return nil
}

func fillSecurities(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListSecurities(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("종목 목록 조회 실패: %w", err)
	}

	snapshot.Securities = make([]SecuritySummary, 0, len(rows))
	for _, row := range rows {
		snapshot.Securities = append(snapshot.Securities, SecuritySummary{
			ID:         row.ID,
			Name:       row.Name,
			Symbol:     row.Symbol,
			Type:       row.Type,
			CategoryID: nullString(row.CategoryID),
			Currency:   row.Currency,
			Note:       row.Note,
		})
	}
	return nil
}

func fillCategoryAllocations(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListCategoryAllocations(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("카테고리 배분 조회 실패: %w", err)
	}

	snapshot.CategoryAllocations = make([]CategoryAllocation, 0, len(rows))
	for _, row := range rows {
		amount, err := toInt64(row.Amount)
		if err != nil {
			return fmt.Errorf("카테고리 금액 변환 실패: %w", err)
		}
		snapshot.CategoryAllocations = append(snapshot.CategoryAllocations, CategoryAllocation{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Amount:       amount,
			WeightPct:    row.WeightPct,
		})
	}
	return nil
}

func fillAccounts(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListAccounts(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("계좌 조회 실패: %w", err)
	}

	snapshot.Accounts = make([]Account, 0, len(rows))
	for _, row := range rows {
		snapshot.Accounts = append(snapshot.Accounts, Account{
			ID:             row.ID,
			Name:           row.Name,
			Provider:       row.Provider,
			CategoryName:   row.CategoryName,
			Balance:        row.Balance,
			MonthlyContrib: row.MonthlyContrib,
			Note:           row.Note,
		})
	}
	return nil
}

func fillHoldings(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListHoldings(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("보유 종목 조회 실패: %w", err)
	}

	snapshot.Holdings = make([]Holding, 0, len(rows))
	for _, row := range rows {
		snapshot.Holdings = append(snapshot.Holdings, Holding{
			ID:             row.ID,
			SecurityID:     row.SecurityID,
			SecurityName:   row.SecurityName,
			SecuritySymbol: row.SecuritySymbol,
			SecurityType:   row.SecurityType,
			CategoryName:   nullString(row.CategoryName),
			Currency:       row.SecurityCurrency,
			Amount:         row.Amount,
			TargetAmount:   row.TargetAmount,
			Note:           row.Note,
		})
	}
	return nil
}

func fillRebalanceGaps(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListRebalanceGaps(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("리밸런싱 현황 조회 실패: %w", err)
	}

	snapshot.RebalanceGaps = make([]RebalanceGap, 0, len(rows))
	for _, row := range rows {
		currentAmount, err := toInt64(row.CurrentAmount)
		if err != nil {
			return fmt.Errorf("현재 금액 변환 실패: %w", err)
		}
		targetFinal, err := toInt64(row.TargetFinal)
		if err != nil {
			return fmt.Errorf("목표 금액 변환 실패: %w", err)
		}
		snapshot.RebalanceGaps = append(snapshot.RebalanceGaps, RebalanceGap{
			CategoryID:     row.CategoryID,
			CategoryName:   row.CategoryName,
			TargetWeight:   row.TargetWeight,
			TargetAmount:   row.TargetAmount,
			CurrentAmount:  currentAmount,
			TargetByWeight: int64(math.Round(row.TargetByWeight)),
			TargetFinal:    targetFinal,
			GapAmount:      row.GapAmount,
		})
	}
	return nil
}

func fillBudgetEntries(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListBudgetEntries(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("예산 항목 조회 실패: %w", err)
	}

	snapshot.BudgetEntries = make([]BudgetEntry, 0, len(rows))
	for _, row := range rows {
		snapshot.BudgetEntries = append(snapshot.BudgetEntries, BudgetEntry{
			ID:        row.ID,
			Name:      row.Name,
			Direction: row.Direction,
			Class:     row.Class,
			Planned:   row.Planned,
			Actual:    row.Actual,
			Note:      row.Note,
		})
	}
	return nil
}

func fillContributionPlans(ctx context.Context, queries *db.Queries, snapshot *PortfolioSnapshot) error {
	rows, err := queries.ListContributionPlans(ctx, snapshot.UID)
	if err != nil {
		return fmt.Errorf("적립 계획 조회 실패: %w", err)
	}

	snapshot.ContributionPlans = make([]ContributionPlan, 0, len(rows))
	for _, row := range rows {
		snapshot.ContributionPlans = append(snapshot.ContributionPlans, ContributionPlan{
			ID:             row.ID,
			SecurityID:     row.SecurityID,
			SecurityName:   row.SecurityName,
			SecuritySymbol: row.SecuritySymbol,
			CategoryName:   nullString(row.CategoryName),
			Weight:         row.Weight,
			Amount:         row.Amount,
			Note:           row.Note,
		})
	}
	return nil
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case float64:
		return int64(math.Round(v)), nil
	case float32:
		return int64(math.Round(float64(v))), nil
	case []byte:
		return parseNumericString(string(v))
	case string:
		return parseNumericString(v)
	default:
		return 0, fmt.Errorf("지원하지 않는 숫자 타입: %T", value)
	}
}

func parseNumericString(raw string) (int64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	if strings.Contains(trimmed, ".") {
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0, err
		}
		return int64(math.Round(f)), nil
	}
	i, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
