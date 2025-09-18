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
			Main(Class("container"),
				Div(Class("max-width-3 margin-auto padding"),
					H1(Class("margin-bottom"), Text("포트폴리오 대시보드")),
					If(snapshot.IsDemo,
						Div(Class("note warning margin-bottom"),
							Text("로그인이 되어 있지 않아 데모 데이터를 보여주고 있어요."),
						),
					),
					Div(
						ID("portfolio-dashboard"),
						DashboardContent(snapshot),
					),
				),
			),
		},
	})
}

func DashboardContent(snapshot *services.PortfolioSnapshot) Node {
	return Group([]Node{
		TotalSection(snapshot),
		CategorySection(snapshot),
		RebalanceSection(snapshot),
		SecuritySection(snapshot),
		AccountSection(snapshot),
		HoldingSection(snapshot),
		BudgetSection(snapshot),
		ContributionSection(snapshot),
	})
}

func cardSection(id string, title string, form Node, headers []string, rows []Node) Node {
	headerNodes := make([]Node, 0, len(headers))
	for _, text := range headers {
		headerNodes = append(headerNodes, Th(Text(text)))
	}

	children := []Node{H2(Text(title))}
	if form != nil {
		children = append(children, Div(Class("margin-bottom"), form))
	}
	children = append(children,
		Table(
			THead(Tr(headerNodes...)),
			TBody(Group(rows)),
		),
	)

	return Section(
		ID(id),
		Class("card margin-bottom"),
		Group(children),
	)
}

func TotalSection(snapshot *services.PortfolioSnapshot) Node {
	return Section(
		ID("total-section"),
		Class("card margin-bottom"),
		H2(Text("총자산")),
		Strong(Class("text-large"), Text(formatCurrency(snapshot.TotalAsset))),
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

	return cardSection(
		"budget-section",
		"월 수입/지출",
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
		Class("grid gap-small"),
		Method("post"),
		Attr("hx-post", action),
		Attr("hx-target", "#portfolio-dashboard"),
		Attr("hx-swap", "innerHTML"),
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
