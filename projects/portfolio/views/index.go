package views

import (
	"fmt"
	"strconv"
	"strings"

	"simple-server/pkg/util/gomutil"
	"simple-server/projects/portfolio/services"
	shared "simple-server/shared/views"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

var (
	categoryRoleOptions = []struct {
		value string
		label string
	}{
		{"liquidity", "유동성"},
		{"growth", "성장"},
		{"income", "인컴"},
		{"protection", "안정"},
		{"other", "기타"},
	}
	securityTypeOptions = []struct {
		value string
		label string
	}{
		{"stock", "주식"},
		{"etf", "ETF"},
		{"fund", "펀드"},
		{"bond", "채권"},
		{"crypto", "가상자산"},
		{"other", "기타"},
	}
	budgetDirectionOptions = []struct {
		value string
		label string
	}{
		{"income", "수입"},
		{"expense", "지출"},
	}
	budgetClassOptions = []struct {
		value string
		label string
	}{
		{"fixed", "고정"},
		{"variable", "유동"},
	}
)

func categorySelectOptions(categories []services.CategorySummary) []Node {
	options := make([]Node, 0, len(categories)+1)
	options = append(options, Option(Value(""), Text("선택해 주세요")))
	for _, category := range categories {
		label := category.Name
		if roleLabel := categoryRoleLabel(category.Role); roleLabel != "" {
			label = fmt.Sprintf("%s (%s)", label, roleLabel)
		}
		options = append(options, Option(Value(category.ID), Text(label)))
	}
	return options
}

func securitySelectOptions(securities []services.SecuritySummary) []Node {
	options := make([]Node, 0, len(securities)+1)
	options = append(options, Option(Value(""), Text("선택해 주세요")))
	for _, security := range securities {
		label := security.Name
		if strings.TrimSpace(security.Symbol) != "" {
			label = fmt.Sprintf("%s (%s)", label, security.Symbol)
		}
		label = fmt.Sprintf("%s - %s", label, securityTypeLabel(security.Type))
		options = append(options, Option(Value(security.ID), Text(label)))
	}
	return options
}

func optionsFromList(items []struct {
	value string
	label string
}) []Node {
	options := []Node{Option(
		Value(""),
		Attr("disabled", "disabled"),
		Attr("selected", "selected"),
		Text("선택해 주세요"),
	)}
	for _, item := range items {
		options = append(options, Option(Value(item.value), Text(item.label)))
	}
	return options
}

func Index(title string, snapshot *services.PortfolioSnapshot) Node {
	if snapshot == nil {
		snapshot = &services.PortfolioSnapshot{}
	}

	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Meta(Name("description"), Content("자산 포트폴리오")),
				Link(Rel("manifest"), Href("/manifest.json")),
			},
		),
		Body: []Node{
			shared.Snackbar(),
			Main(Class("portfolio-page"),
				Div(Class("portfolio-shell"),
					Div(
						ID("portfolio-dashboard"),
						Class("portfolio-dashboard"),
						DashboardContent(snapshot),
					),
				),
			),
		},
	})
}

func DashboardContent(snapshot *services.PortfolioSnapshot) Node {
	return Group([]Node{
		OverviewSection(snapshot),
		Div(
			Class("grid large-space responsive portfolio-grid"),
			Div(Class("s12 l6"), CategorySection(snapshot)),
			Div(Class("s12 l6"), RebalanceSection(snapshot)),
			Div(Class("s12 l6"), AccountSection(snapshot)),
			Div(Class("s12 l6"), HoldingSection(snapshot)),
			Div(Class("s12 l6"), SecuritySection(snapshot)),
			Div(Class("s12 l6"), ContributionSection(snapshot)),
			Div(Class("s12"), BudgetSection(snapshot)),
		),
	})
}

func OverviewSection(snapshot *services.PortfolioSnapshot) Node {
	monthlyContribution := sumMonthlyContribution(snapshot.Accounts)
	monthlyPlan := sumContributionAmount(snapshot.ContributionPlans)
	incomePlan, expensePlan, incomeActual, expenseActual := budgetTotals(snapshot.BudgetEntries)
	plannedNet := incomePlan - expensePlan
	actualNet := incomeActual - expenseActual

	headerColumns := []Node{
		Div(
			Class("s12 m8 l7 hero-text"),
			H1(Class("hero-title"), Text("포트폴리오 대시보드")),
			P(Class("hero-description"),
				Text("현재 자산 현황과 월간 계획을 한눈에 확인하고 필요한 항목을 바로 추가할 수 있어요."),
			),
		),
	}

	if snapshot.IsDemo {
		headerColumns = append(headerColumns,
			Div(
				Class("s12 m4 l5 hero-actions"),
				Div(Class("chip demo-chip"),
					I(Class("icon"), Text("visibility")),
					Span(Text("로그인하지 않아 데모 데이터를 보고 있어요")),
				),
			),
		)
	}

	return Section(
		Class("portfolio-hero surface-container-highest shadow round"),
		Div(
			Class("hero-headline grid large-space responsive"),
			Group(headerColumns),
		),
		Div(
			Class("hero-metrics grid large-space responsive"),
			Div(Class("s12 m6 l3"), overviewMetric("총자산", formatCurrency(snapshot.TotalAsset), fmt.Sprintf("카테고리 %d개", len(snapshot.Categories)))),
			Div(Class("s12 m6 l3"), overviewMetric("계좌", fmt.Sprintf("%d개", len(snapshot.Accounts)), fmt.Sprintf("월 납입 %s", formatCurrency(monthlyContribution)))),
			Div(Class("s12 m6 l3"), overviewMetric("보유 종목", fmt.Sprintf("%d개", len(snapshot.Holdings)), fmt.Sprintf("평가금액 %s", formatCurrency(sumHoldingsAmount(snapshot.Holdings))))),
			Div(Class("s12 m6 l3"), overviewMetric("월 적립 계획", formatCurrency(monthlyPlan), fmt.Sprintf("실행 중 %d건", len(snapshot.ContributionPlans)))),
		),
		Div(
			Class("hero-budget grid responsive right-align"),
			Div(Class("budget-card surface-container-high shadow round s12 m8 l4"),
				Div(Class("budget-label"), Text("월 수입/지출 요약")),
				Div(Class("budget-values"),
					Div(Class("budget-column"),
						Span(Class("budget-chip chip"), Text("계획")),
						Strong(Class("budget-amount"), Text(formatSignedCurrency(plannedNet))),
					),
					Div(Class("budget-column"),
						Span(Class("budget-chip chip"), Text("실적")),
						Strong(Class("budget-amount"), Text(formatSignedCurrency(actualNet))),
					),
				),
				P(Class("budget-helper"),
					Text("양수는 잉여, 음수는 부족을 의미해요."),
				),
			),
		),
	)
}

func overviewMetric(title string, value string, helper string) Node {
	return Div(Class("metric-card surface-container-high shadow round"),
		Span(Class("metric-title"), Text(title)),
		Strong(Class("metric-value"), Text(value)),
		If(helper != "",
			Span(Class("metric-helper"), Text(helper)),
		),
	)
}

func cardSection(id string, title string, subtitle string, form Node, headers []string, rows []Node) Node {
	headerNodes := make([]Node, 0, len(headers))
	for _, text := range headers {
		headerNodes = append(headerNodes, Th(Text(text)))
	}

	headerChildren := []Node{
		Div(Class("card-heading"),
			H2(Class("card-title"), Text(title)),
			If(subtitle != "",
				P(Class("card-subtitle"), Text(subtitle)),
			),
		),
	}

	if panel := ifFormPanel(title, form); panel != nil {
		headerChildren = append(headerChildren, panel)
	}

	return Section(
		ID(id),
		Class("portfolio-card surface-container-high shadow round"),
		Group([]Node{
			Div(Class("card-header"), Group(headerChildren)),
			Div(Class("table-container scroll"),
				Table(
					THead(Tr(headerNodes...)),
					TBody(Group(rows)),
				),
			),
		}),
	)
}

func CategorySection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"카테고리", "현재 금액", "비중"}
	rows := make([]Node, 0, len(snapshot.CategoryAllocations))
	for _, item := range snapshot.CategoryAllocations {
		rows = append(rows, Tr(
			Td(Text(item.CategoryName)),
			Td(Text(formatCurrency(item.Amount))),
			Td(Text(formatPercent(item.WeightPct))),
		))
	}
	rows = ensureRows(rows, len(headers), "카테고리 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = categoryForm(snapshot.Categories)
	}

	return cardSection(
		"category-section",
		"카테고리별 자산 배분",
		fmt.Sprintf("총 %s · %d개 카테고리", formatCurrency(sumCategoryAmount(snapshot.CategoryAllocations)), len(snapshot.Categories)),
		form,
		headers,
		rows,
	)
}

func RebalanceSection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"카테고리", "목표 비중", "목표 금액", "비중 기준 목표", "현재 금액", "차이"}
	rows := make([]Node, 0, len(snapshot.RebalanceGaps))
	for _, item := range snapshot.RebalanceGaps {
		rows = append(rows, Tr(
			Td(Text(item.CategoryName)),
			Td(Text(formatPercent(item.TargetWeight))),
			Td(Text(formatCurrency(item.TargetAmount))),
			Td(Text(formatCurrency(item.TargetByWeight))),
			Td(Text(formatCurrency(item.CurrentAmount))),
			Td(Text(formatSignedCurrency(item.GapAmount))),
		))
	}
	rows = ensureRows(rows, len(headers), "리밸런싱 목표 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = allocationTargetForm(snapshot.Categories)
	}

	return cardSection(
		"rebalance-section",
		"리밸런싱 목표 대비 현황",
		fmt.Sprintf("목표 %d개 · 추가 필요 %s", len(snapshot.RebalanceGaps), formatCurrency(sumPositiveGap(snapshot.RebalanceGaps))),
		form,
		headers,
		rows,
	)
}

func SecuritySection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"종목", "심볼", "유형", "카테고리", "통화", "메모"}
	rows := make([]Node, 0, len(snapshot.Securities))
	for _, item := range snapshot.Securities {
		rows = append(rows, Tr(
			Td(Text(item.Name)),
			Td(Text(emptyFallback(item.Symbol))),
			Td(Text(securityTypeLabel(item.Type))),
			Td(Text(emptyFallback(findCategoryName(snapshot.Categories, item.CategoryID)))),
			Td(Text(item.Currency)),
			Td(Text(formatNote(item.Note))),
		))
	}
	rows = ensureRows(rows, len(headers), "종목 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = securityForm(snapshot.Categories)
	}

	return cardSection(
		"security-section",
		"종목 마스터",
		fmt.Sprintf("등록 %d개 · 카테고리 연결 %d개", len(snapshot.Securities), countLinkedSecurities(snapshot.Securities)),
		form,
		headers,
		rows,
	)
}

func AccountSection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"계좌명", "기관", "카테고리", "잔액", "월 납입", "메모"}
	rows := make([]Node, 0, len(snapshot.Accounts))
	for _, item := range snapshot.Accounts {
		rows = append(rows, Tr(
			Td(Text(item.Name)),
			Td(Text(item.Provider)),
			Td(Text(item.CategoryName)),
			Td(Text(formatCurrency(item.Balance))),
			Td(Text(formatCurrency(item.MonthlyContrib))),
			Td(Text(formatNote(item.Note))),
		))
	}
	rows = ensureRows(rows, len(headers), "계좌 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = accountForm(snapshot.Categories)
	}

	return cardSection(
		"account-section",
		"계좌 현황",
		fmt.Sprintf("총 잔액 %s · 월 납입 %s", formatCurrency(sumAccountBalance(snapshot.Accounts)), formatCurrency(sumMonthlyContribution(snapshot.Accounts))),
		form,
		headers,
		rows,
	)
}

func HoldingSection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"종목", "심볼", "유형", "카테고리", "현재 금액", "목표 금액", "메모"}
	rows := make([]Node, 0, len(snapshot.Holdings))
	for _, item := range snapshot.Holdings {
		rows = append(rows, Tr(
			Td(Text(item.SecurityName)),
			Td(Text(item.SecuritySymbol)),
			Td(Text(securityTypeLabel(item.SecurityType))),
			Td(Text(emptyFallback(item.CategoryName))),
			Td(Text(formatCurrency(item.Amount))),
			Td(Text(formatCurrency(item.TargetAmount))),
			Td(Text(formatNote(item.Note))),
		))
	}
	rows = ensureRows(rows, len(headers), "보유 종목 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = holdingForm(snapshot.Securities)
	}

	return cardSection(
		"holding-section",
		"보유 종목",
		fmt.Sprintf("평가금액 %s · 목표 합계 %s (%d건)",
			formatCurrency(sumHoldingsAmount(snapshot.Holdings)),
			formatCurrency(sumHoldingTarget(snapshot.Holdings)),
			countHoldingsWithTarget(snapshot.Holdings),
		),
		form,
		headers,
		rows,
	)
}

func BudgetSection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"항목", "구분", "유형", "계획", "실적", "메모"}
	rows := make([]Node, 0, len(snapshot.BudgetEntries))
	for _, item := range snapshot.BudgetEntries {
		rows = append(rows, Tr(
			Td(Text(item.Name)),
			Td(Text(directionLabel(item.Direction))),
			Td(Text(classLabel(item.Class))),
			Td(Text(formatCurrency(item.Planned))),
			Td(Text(formatCurrency(item.Actual))),
			Td(Text(formatNote(item.Note))),
		))
	}
	rows = ensureRows(rows, len(headers), "예산 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = budgetForm()
	}

	incomePlan, expensePlan, incomeActual, expenseActual := budgetTotals(snapshot.BudgetEntries)

	return cardSection(
		"budget-section",
		"월 수입/지출",
		fmt.Sprintf("계획 수입 %s · 지출 %s | 실적 %s · %s",
			formatCurrency(incomePlan),
			formatCurrency(expensePlan),
			formatCurrency(incomeActual),
			formatCurrency(expenseActual),
		),
		form,
		headers,
		rows,
	)
}

func ContributionSection(snapshot *services.PortfolioSnapshot) Node {
	headers := []string{"종목", "심볼", "카테고리", "비중", "월 적립액", "메모"}
	rows := make([]Node, 0, len(snapshot.ContributionPlans))
	for _, item := range snapshot.ContributionPlans {
		rows = append(rows, Tr(
			Td(Text(item.SecurityName)),
			Td(Text(item.SecuritySymbol)),
			Td(Text(emptyFallback(item.CategoryName))),
			Td(Text(formatPercentOne(item.Weight))),
			Td(Text(formatCurrency(item.Amount))),
			Td(Text(formatNote(item.Note))),
		))
	}
	rows = ensureRows(rows, len(headers), "월 적립 계획 데이터가 없습니다.")

	var form Node
	if !snapshot.IsDemo {
		form = contributionForm(snapshot.Securities)
	}

	return cardSection(
		"contribution-section",
		"주식 월 적립 계획",
		fmt.Sprintf("월 적립 합계 %s · %d건", formatCurrency(sumContributionAmount(snapshot.ContributionPlans)), len(snapshot.ContributionPlans)),
		form,
		headers,
		rows,
	)
}

func hxForm(action string, submitLabel string, fields ...Node) Node {
	nodes := make([]Node, 0, len(fields)+1)
	nodes = append(nodes, fields...)
	nodes = append(nodes, Div(Class("form-actions"), Button(Class("button primary"), Type("submit"), Text(submitLabel))))

	return Form(
		Class("form-grid"),
		Method("post"),
		Attr("hx-post", action),
		Attr("hx-target", "#portfolio-dashboard"),
		Attr("hx-swap", "innerHTML"),
		Attr("hx-on::after-request", "if (event.detail.successful) this.reset()"),
		Group(nodes),
	)
}

func categoryForm(categories []services.CategorySummary) Node {
	return hxForm(
		"/categories",
		"카테고리 추가",
		Div(Class("field"),
			Label(For("category-name"), Text("카테고리명")),
			Input(ID("category-name"), Name("name"), Type("text"), Attr("required", "required")),
		),
		Div(Class("field"),
			Label(For("category-role"), Text("역할")),
			Select(
				ID("category-role"),
				Name("role"),
				Attr("required", "required"),
				Group(optionsFromList(categoryRoleOptions)),
			),
		),
		Div(Class("field"),
			Label(For("category-parent"), Text("상위 카테고리")),
			Select(
				ID("category-parent"),
				Name("parent_id"),
				Group(categorySelectOptions(categories)),
			),
		),
		Div(Class("field"),
			Label(For("category-order"), Text("정렬 순서")),
			Input(ID("category-order"), Name("display_order"), Type("number"), Attr("min", "0")),
		),
	)
}

func ifFormPanel(title string, form Node) Node {
	if form == nil {
		return nil
	}

	label := "새 항목 추가"
	if parts := strings.Fields(title); len(parts) > 0 {
		label = fmt.Sprintf("%s 추가", parts[0])
	}

	return Details(
		Class("form-panel"),
		Summary(
			Class("form-summary"),
			I(Class("icon"), Text("add")),
			Span(Text(label)),
		),
		Div(Class("form-body"), form),
	)
}

func allocationTargetForm(categories []services.CategorySummary) Node {
	return hxForm(
		"/allocation-targets",
		"목표 추가",
		Div(Class("field"),
			Label(For("target-category"), Text("카테고리")),
			Select(
				ID("target-category"),
				Name("category_id"),
				Attr("required", "required"),
				Group(categorySelectOptions(categories)),
			),
		),
		Div(Class("field"),
			Label(For("target-weight"), Text("목표 비중(%)")),
			Input(ID("target-weight"), Name("target_weight"), Type("number"), Attr("step", "0.1"), Attr("min", "0"), Attr("max", "100")),
		),
		Div(Class("field"),
			Label(For("target-amount"), Text("목표 금액")),
			Input(ID("target-amount"), Name("target_amount"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("target-note"), Text("메모")),
			Textarea(ID("target-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func securityForm(categories []services.CategorySummary) Node {
	return hxForm(
		"/securities",
		"종목 추가",
		Div(Class("field"),
			Label(For("security-name"), Text("종목명")),
			Input(ID("security-name"), Name("name"), Type("text"), Attr("required", "required")),
		),
		Div(Class("field"),
			Label(For("security-symbol"), Text("심볼")),
			Input(ID("security-symbol"), Name("symbol"), Type("text")),
		),
		Div(Class("field"),
			Label(For("security-type"), Text("유형")),
			Select(
				ID("security-type"),
				Name("type"),
				Attr("required", "required"),
				Group(optionsFromList(securityTypeOptions)),
			),
		),
		Div(Class("field"),
			Label(For("security-category"), Text("카테고리")),
			Select(
				ID("security-category"),
				Name("category_id"),
				Group(categorySelectOptions(categories)),
			),
		),
		Div(Class("field"),
			Label(For("security-currency"), Text("통화")),
			Input(ID("security-currency"), Name("currency"), Type("text"), Attr("value", "KRW")),
		),
		Div(Class("field"),
			Label(For("security-note"), Text("메모")),
			Textarea(ID("security-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func accountForm(categories []services.CategorySummary) Node {
	return hxForm(
		"/accounts",
		"계좌 추가",
		Div(Class("field"),
			Label(For("account-name"), Text("계좌명")),
			Input(ID("account-name"), Name("name"), Type("text"), Attr("required", "required")),
		),
		Div(Class("field"),
			Label(For("account-provider"), Text("기관")),
			Input(ID("account-provider"), Name("provider"), Type("text")),
		),
		Div(Class("field"),
			Label(For("account-category"), Text("카테고리")),
			Select(
				ID("account-category"),
				Name("category_id"),
				Attr("required", "required"),
				Group(categorySelectOptions(categories)),
			),
		),
		Div(Class("field"),
			Label(For("account-balance"), Text("잔액")),
			Input(ID("account-balance"), Name("balance"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("account-monthly"), Text("월 납입")),
			Input(ID("account-monthly"), Name("monthly_contrib"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("account-note"), Text("메모")),
			Textarea(ID("account-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func holdingForm(securities []services.SecuritySummary) Node {
	return hxForm(
		"/holdings",
		"보유 종목 추가",
		Div(Class("field"),
			Label(For("holding-security"), Text("종목")),
			Select(
				ID("holding-security"),
				Name("security_id"),
				Attr("required", "required"),
				Group(securitySelectOptions(securities)),
			),
		),
		Div(Class("field"),
			Label(For("holding-amount"), Text("현재 금액")),
			Input(ID("holding-amount"), Name("amount"), Type("number"), Attr("min", "0"), Attr("required", "required")),
		),
		Div(Class("field"),
			Label(For("holding-target"), Text("목표 금액")),
			Input(ID("holding-target"), Name("target_amount"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("holding-note"), Text("메모")),
			Textarea(ID("holding-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func budgetForm() Node {
	return hxForm(
		"/budget-entries",
		"예산 추가",
		Div(Class("field"),
			Label(For("budget-name"), Text("항목")),
			Input(ID("budget-name"), Name("name"), Type("text"), Attr("required", "required")),
		),
		Div(Class("field"),
			Label(For("budget-direction"), Text("구분")),
			Select(
				ID("budget-direction"),
				Name("direction"),
				Attr("required", "required"),
				Group(optionsFromList(budgetDirectionOptions)),
			),
		),
		Div(Class("field"),
			Label(For("budget-class"), Text("유형")),
			Select(
				ID("budget-class"),
				Name("class"),
				Attr("required", "required"),
				Group(optionsFromList(budgetClassOptions)),
			),
		),
		Div(Class("field"),
			Label(For("budget-planned"), Text("계획 금액")),
			Input(ID("budget-planned"), Name("planned"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("budget-actual"), Text("실제 금액")),
			Input(ID("budget-actual"), Name("actual"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("budget-note"), Text("메모")),
			Textarea(ID("budget-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func contributionForm(securities []services.SecuritySummary) Node {
	return hxForm(
		"/contribution-plans",
		"적립 계획 추가",
		Div(Class("field"),
			Label(For("plan-security"), Text("종목")),
			Select(
				ID("plan-security"),
				Name("security_id"),
				Attr("required", "required"),
				Group(securitySelectOptions(securities)),
			),
		),
		Div(Class("field"),
			Label(For("plan-weight"), Text("비중(%)")),
			Input(ID("plan-weight"), Name("weight"), Type("number"), Attr("step", "0.1"), Attr("min", "0"), Attr("max", "100")),
		),
		Div(Class("field"),
			Label(For("plan-amount"), Text("월 적립액")),
			Input(ID("plan-amount"), Name("amount"), Type("number"), Attr("min", "0")),
		),
		Div(Class("field"),
			Label(For("plan-note"), Text("메모")),
			Textarea(ID("plan-note"), Name("note"), Attr("rows", "2")),
		),
	)
}

func ensureRows(rows []Node, cols int, emptyMessage string) []Node {
	if len(rows) > 0 {
		return rows
	}
	return []Node{
		Tr(
			Td(Attr("colspan", strconv.Itoa(cols)), Class("text-center"), Text(emptyMessage)),
		),
	}
}

func formatCurrency(amount int64) string {
	return formatCurrencyInternal(amount, false)
}

func formatSignedCurrency(amount int64) string {
	return formatCurrencyInternal(amount, true)
}

func formatCurrencyInternal(amount int64, signed bool) string {
	sign := ""
	value := amount
	if value < 0 {
		sign = "-"
		value = -value
	} else if signed {
		sign = "+"
	}

	digits := strconv.FormatInt(value, 10)
	var builder strings.Builder
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune(digit)
	}

	return fmt.Sprintf("%s%s원", sign, builder.String())
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func formatPercentOne(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

func formatNote(note string) string {
	if strings.TrimSpace(note) == "" {
		return "-"
	}
	return note
}

func categoryRoleLabel(role string) string {
	for _, item := range categoryRoleOptions {
		if item.value == role {
			return item.label
		}
	}
	return ""
}

func findCategoryName(categories []services.CategorySummary, id string) string {
	for _, category := range categories {
		if category.ID == id {
			return category.Name
		}
	}
	return ""
}

func directionLabel(direction string) string {
	switch direction {
	case "income":
		return "수입"
	case "expense":
		return "지출"
	default:
		return direction
	}
}

func classLabel(class string) string {
	switch class {
	case "fixed":
		return "고정"
	case "variable":
		return "유동"
	default:
		return class
	}
}

func securityTypeLabel(t string) string {
	switch t {
	case "stock":
		return "주식"
	case "etf":
		return "ETF"
	case "fund":
		return "펀드"
	case "bond":
		return "채권"
	case "crypto":
		return "가상자산"
	default:
		return "기타"
	}
}

func emptyFallback(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func sumCategoryAmount(items []services.CategoryAllocation) int64 {
	var total int64
	for _, item := range items {
		total += item.Amount
	}
	return total
}

func sumAccountBalance(accounts []services.Account) int64 {
	var total int64
	for _, account := range accounts {
		total += account.Balance
	}
	return total
}

func sumMonthlyContribution(accounts []services.Account) int64 {
	var total int64
	for _, account := range accounts {
		total += account.MonthlyContrib
	}
	return total
}

func sumHoldingsAmount(holdings []services.Holding) int64 {
	var total int64
	for _, holding := range holdings {
		total += holding.Amount
	}
	return total
}

func sumHoldingTarget(holdings []services.Holding) int64 {
	var total int64
	for _, holding := range holdings {
		total += holding.TargetAmount
	}
	return total
}

func sumContributionAmount(plans []services.ContributionPlan) int64 {
	var total int64
	for _, plan := range plans {
		total += plan.Amount
	}
	return total
}

func countLinkedSecurities(securities []services.SecuritySummary) int {
	count := 0
	for _, security := range securities {
		if strings.TrimSpace(security.CategoryID) != "" {
			count++
		}
	}
	return count
}

func countHoldingsWithTarget(holdings []services.Holding) int {
	count := 0
	for _, holding := range holdings {
		if holding.TargetAmount > 0 {
			count++
		}
	}
	return count
}

func sumPositiveGap(gaps []services.RebalanceGap) int64 {
	var total int64
	for _, gap := range gaps {
		if gap.GapAmount > 0 {
			total += gap.GapAmount
		}
	}
	return total
}

func budgetTotals(entries []services.BudgetEntry) (int64, int64, int64, int64) {
	var incomePlan int64
	var expensePlan int64
	var incomeActual int64
	var expenseActual int64
	for _, entry := range entries {
		switch entry.Direction {
		case "income":
			incomePlan += entry.Planned
			incomeActual += entry.Actual
		case "expense":
			expensePlan += entry.Planned
			expenseActual += entry.Actual
		}
	}
	return incomePlan, expensePlan, incomeActual, expenseActual
}
