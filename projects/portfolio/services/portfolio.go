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

	if uid != DemoUID {
		if err := ensurePortfolioSeed(ctx, queries, uid); err != nil {
			return nil, fmt.Errorf("포트폴리오 초기 데이터 준비 실패: %w", err)
		}
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

func ensurePortfolioSeed(ctx context.Context, queries *db.Queries, uid string) error {
	count, err := queries.CountCategories(ctx, uid)
	if err != nil {
		return fmt.Errorf("초기 데이터 확인 실패: %w", err)
	}
	if count > 0 {
		return nil
	}

	return cloneDemoPortfolio(ctx, queries, uid)
}

type demoSeed struct {
	Categories []db.Category
	Securities []db.Security
	Accounts   []db.ListAccountsRow
	Holdings   []db.ListHoldingsRow
	Targets    []db.ListAllocationTargetsRow
	Budgets    []db.ListBudgetEntriesRow
	Plans      []db.ListContributionPlansRow
}

func loadDemoSeed(ctx context.Context, queries *db.Queries) (*demoSeed, error) {
	categories, err := queries.ListCategories(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 카테고리 조회 실패: %w", err)
	}
	securities, err := queries.ListSecurities(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 종목 조회 실패: %w", err)
	}
	accounts, err := queries.ListAccounts(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 계좌 조회 실패: %w", err)
	}
	holdings, err := queries.ListHoldings(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 보유 종목 조회 실패: %w", err)
	}
	targets, err := queries.ListAllocationTargets(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 리밸런싱 목표 조회 실패: %w", err)
	}
	budgets, err := queries.ListBudgetEntries(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 예산 항목 조회 실패: %w", err)
	}
	plans, err := queries.ListContributionPlans(ctx, DemoUID)
	if err != nil {
		return nil, fmt.Errorf("데모 적립 계획 조회 실패: %w", err)
	}

	return &demoSeed{
		Categories: categories,
		Securities: securities,
		Accounts:   accounts,
		Holdings:   holdings,
		Targets:    targets,
		Budgets:    budgets,
		Plans:      plans,
	}, nil
}

func cloneDemoPortfolio(ctx context.Context, queries *db.Queries, uid string) error {
	seed, err := loadDemoSeed(ctx, queries)
	if err != nil {
		return err
	}

	conn, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %w", err)
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("데모 데이터 복제 트랜잭션 시작 실패: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var (
		categoryMap map[string]string
		securityMap map[string]string
	)

	steps := []func() error{
		func() error {
			var err error
			categoryMap, err = copyCategories(ctx, tx, seed.Categories, uid)
			return err
		},
		func() error {
			var err error
			securityMap, err = copySecurities(ctx, tx, seed.Securities, uid, categoryMap)
			return err
		},
		func() error {
			return copyAccounts(ctx, tx, seed.Accounts, uid, categoryMap)
		},
		func() error {
			return copyHoldings(ctx, tx, seed.Holdings, uid, securityMap)
		},
		func() error {
			return copyAllocationTargets(ctx, tx, seed.Targets, uid, categoryMap)
		},
		func() error {
			return copyBudgetEntries(ctx, tx, seed.Budgets, uid)
		},
		func() error {
			return copyContributionPlans(ctx, tx, seed.Plans, uid, securityMap)
		},
	}

	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("데모 데이터 복제 커밋 실패: %w", err)
	}
	return nil
}

func copyCategories(ctx context.Context, tx *sql.Tx, categories []db.Category, uid string) (map[string]string, error) {
	mapping := make(map[string]string, len(categories))
	remaining := append([]db.Category(nil), categories...)

	for len(remaining) > 0 {
		progressed := false
		next := make([]db.Category, 0)

		for _, cat := range remaining {
			var parent interface{}
			if cat.ParentID.Valid {
				newParent, ok := mapping[cat.ParentID.String]
				if !ok {
					next = append(next, cat)
					continue
				}
				parent = newParent
			}

			newID, err := insertCategory(ctx, tx, uid, cat.Name, cat.Role, parent, cat.DisplayOrder)
			if err != nil {
				return nil, fmt.Errorf("카테고리 복제 실패(%s): %w", cat.Name, err)
			}
			mapping[cat.ID] = newID
			progressed = true
		}

		if !progressed {
			return nil, fmt.Errorf("카테고리 복제 실패: 상위 카테고리를 찾을 수 없습니다")
		}

		remaining = next
	}

	return mapping, nil
}

func copySecurities(ctx context.Context, tx *sql.Tx, securities []db.Security, uid string, categoryMap map[string]string) (map[string]string, error) {
	mapping := make(map[string]string, len(securities))

	for _, sec := range securities {
		var category interface{}
		if sec.CategoryID.Valid {
			newID, ok := categoryMap[sec.CategoryID.String]
			if !ok {
				return nil, fmt.Errorf("종목 복제 실패(%s): 카테고리 매핑을 찾을 수 없습니다", sec.Name)
			}
			category = newID
		}

		newID, err := insertSecurity(ctx, tx, uid, sec.Symbol, sec.Name, sec.Type, category, sec.Currency, sec.Note)
		if err != nil {
			return nil, fmt.Errorf("종목 복제 실패(%s): %w", sec.Name, err)
		}
		mapping[sec.ID] = newID
	}

	return mapping, nil
}

func copyAccounts(ctx context.Context, tx *sql.Tx, accounts []db.ListAccountsRow, uid string, categoryMap map[string]string) error {
	for _, acc := range accounts {
		newCategory, ok := categoryMap[acc.CategoryID]
		if !ok {
			return fmt.Errorf("계좌 복제 실패(%s): 카테고리 매핑을 찾을 수 없습니다", acc.Name)
		}

		if _, err := tx.ExecContext(ctx, `
                        INSERT INTO account (uid, category_id, name, provider, balance, monthly_contrib, note)
                        VALUES (?, ?, ?, ?, ?, ?, ?)
                `, uid, newCategory, acc.Name, acc.Provider, acc.Balance, acc.MonthlyContrib, acc.Note); err != nil {
			return fmt.Errorf("계좌 복제 실패(%s): %w", acc.Name, err)
		}
	}
	return nil
}

func copyHoldings(ctx context.Context, tx *sql.Tx, holdings []db.ListHoldingsRow, uid string, securityMap map[string]string) error {
	for _, h := range holdings {
		newSecurity, ok := securityMap[h.SecurityID]
		if !ok {
			return fmt.Errorf("보유 종목 복제 실패(%s): 종목 매핑을 찾을 수 없습니다", h.SecurityName)
		}

		if _, err := tx.ExecContext(ctx, `
                        INSERT INTO holding (uid, security_id, amount, target_amount, note)
                        VALUES (?, ?, ?, ?, ?)
                `, uid, newSecurity, h.Amount, h.TargetAmount, h.Note); err != nil {
			return fmt.Errorf("보유 종목 복제 실패(%s): %w", h.SecurityName, err)
		}
	}
	return nil
}

func copyAllocationTargets(ctx context.Context, tx *sql.Tx, targets []db.ListAllocationTargetsRow, uid string, categoryMap map[string]string) error {
	for _, target := range targets {
		newCategory, ok := categoryMap[target.CategoryID]
		if !ok {
			return fmt.Errorf("리밸런싱 목표 복제 실패(%s): 카테고리 매핑을 찾을 수 없습니다", target.CategoryName)
		}

		if _, err := tx.ExecContext(ctx, `
                        INSERT INTO allocation_target (uid, category_id, target_weight, target_amount, note)
                        VALUES (?, ?, ?, ?, ?)
                `, uid, newCategory, target.TargetWeight, target.TargetAmount, target.Note); err != nil {
			return fmt.Errorf("리밸런싱 목표 복제 실패(%s): %w", target.CategoryName, err)
		}
	}
	return nil
}

func copyBudgetEntries(ctx context.Context, tx *sql.Tx, entries []db.ListBudgetEntriesRow, uid string) error {
	for _, entry := range entries {
		if _, err := tx.ExecContext(ctx, `
                        INSERT INTO budget_entry (uid, name, direction, class, planned, actual, note)
                        VALUES (?, ?, ?, ?, ?, ?, ?)
                `, uid, entry.Name, entry.Direction, entry.Class, entry.Planned, entry.Actual, entry.Note); err != nil {
			return fmt.Errorf("예산 항목 복제 실패(%s): %w", entry.Name, err)
		}
	}
	return nil
}

func copyContributionPlans(ctx context.Context, tx *sql.Tx, plans []db.ListContributionPlansRow, uid string, securityMap map[string]string) error {
	for _, plan := range plans {
		newSecurity, ok := securityMap[plan.SecurityID]
		if !ok {
			return fmt.Errorf("적립 계획 복제 실패(%s): 종목 매핑을 찾을 수 없습니다", plan.SecurityName)
		}

		if _, err := tx.ExecContext(ctx, `
                        INSERT INTO contribution_plan (uid, security_id, weight, amount, note)
                        VALUES (?, ?, ?, ?, ?)
                `, uid, newSecurity, plan.Weight, plan.Amount, plan.Note); err != nil {
			return fmt.Errorf("적립 계획 복제 실패(%s): %w", plan.SecurityName, err)
		}
	}
	return nil
}

func insertCategory(ctx context.Context, tx *sql.Tx, uid string, name string, role string, parent interface{}, displayOrder int64) (string, error) {
	if _, err := tx.ExecContext(ctx, `
                INSERT INTO category (uid, name, role, parent_id, display_order)
                VALUES (?, ?, ?, ?, ?)
        `, uid, name, role, parent, displayOrder); err != nil {
		return "", err
	}

	var id string
	if err := tx.QueryRowContext(ctx, `SELECT id FROM category WHERE rowid = last_insert_rowid()`).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func insertSecurity(ctx context.Context, tx *sql.Tx, uid string, symbol string, name string, secType string, category interface{}, currency string, note string) (string, error) {
	if _, err := tx.ExecContext(ctx, `
                INSERT INTO security (uid, symbol, name, type, category_id, currency, note)
                VALUES (?, ?, ?, ?, ?, ?, ?)
        `, uid, symbol, name, secType, category, currency, note); err != nil {
		return "", err
	}

	var id string
	if err := tx.QueryRowContext(ctx, `SELECT id FROM security WHERE rowid = last_insert_rowid()`).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
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
