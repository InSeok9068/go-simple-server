-- sqlc generate -f ./projects/portfolio/sqlc.yaml

-- name: GetUser :one
SELECT * FROM user WHERE uid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO user (uid, name, email) VALUES (?, ?, ?);

-- name: GetTotalAsset :one
SELECT uid, total_amount
FROM v_total_asset
WHERE uid = ?;

-- name: ListCategoryAllocations :many
SELECT v.uid, v.category_id, v.category_name, v.amount, v.weight_pct
FROM v_category_allocation v
JOIN category c ON c.id = v.category_id AND c.uid = v.uid
WHERE v.uid = ?
ORDER BY c.display_order;

-- name: ListAccounts :many
SELECT
    a.id,
    a.uid,
    a.category_id,
    a.name,
    a.provider,
    a.balance,
    a.monthly_contrib,
    a.note,
    c.name AS category_name
FROM account a
JOIN category c ON c.id = a.category_id
WHERE a.uid = ?
ORDER BY c.display_order, a.name;

-- name: ListHoldings :many
SELECT
    h.id,
    h.uid,
    h.security_id,
    h.amount,
    h.target_amount,
    h.note,
    s.symbol AS security_symbol,
    s.name AS security_name,
    s.type AS security_type,
    s.category_id AS security_category_id,
    c.name AS category_name,
    s.currency AS security_currency
FROM holding h
JOIN security s ON s.id = h.security_id
LEFT JOIN category c ON c.id = s.category_id
WHERE h.uid = ?
ORDER BY c.display_order, s.name;

-- name: ListAllocationTargets :many
SELECT
    a.id,
    a.uid,
    a.category_id,
    a.target_weight,
    a.target_amount,
    a.note,
    c.name AS category_name
FROM allocation_target a
JOIN category c ON c.id = a.category_id
WHERE a.uid = ?
ORDER BY c.display_order;

-- name: ListRebalanceGaps :many
SELECT
    g.uid,
    g.category_id,
    c.name AS category_name,
    g.target_weight,
    g.target_amount,
    g.current_amount,
    g.target_by_weight,
    g.target_final,
    g.gap_amount
FROM v_rebalance_gap g
JOIN category c ON c.id = g.category_id
WHERE g.uid = ?
ORDER BY c.display_order;

-- name: ListBudgetEntries :many
SELECT
    id,
    uid,
    name,
    direction,
    class,
    planned,
    actual,
    note
FROM budget_entry
WHERE uid = ?
ORDER BY
    CASE direction WHEN 'income' THEN 0 ELSE 1 END,
    CASE class WHEN 'fixed' THEN 0 ELSE 1 END,
    name;

-- name: ListContributionPlans :many
SELECT
    cp.id,
    cp.uid,
    cp.security_id,
    cp.weight,
    cp.amount,
    cp.note,
    s.symbol AS security_symbol,
    s.name AS security_name,
    s.category_id AS security_category_id,
    c.name AS category_name
FROM contribution_plan cp
JOIN security s ON s.id = cp.security_id
LEFT JOIN category c ON c.id = s.category_id
WHERE cp.uid = ?
ORDER BY cp.weight DESC, s.name;

-- name: ListCategories :many
SELECT
    id,
    uid,
    name,
    role,
    parent_id,
    display_order
FROM category
WHERE uid = ?
ORDER BY display_order, name;

-- name: ListSecurities :many
SELECT
    id,
    uid,
    symbol,
    name,
    type,
    category_id,
    currency,
    note
FROM security
WHERE uid = ?
ORDER BY name;

-- name: CreateCategory :exec
INSERT INTO category (uid, name, role, parent_id, display_order)
VALUES (?, ?, ?, ?, ?);

-- name: CreateSecurity :exec
INSERT INTO security (uid, symbol, name, type, category_id, currency, note)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: CreateAccount :exec
INSERT INTO account (uid, category_id, name, provider, balance, monthly_contrib, note)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: CreateHolding :exec
INSERT INTO holding (uid, security_id, amount, target_amount, note)
VALUES (?, ?, ?, ?, ?);

-- name: CreateAllocationTarget :exec
INSERT INTO allocation_target (uid, category_id, target_weight, target_amount, note)
VALUES (?, ?, ?, ?, ?);

-- name: CreateBudgetEntry :exec
INSERT INTO budget_entry (uid, name, direction, class, planned, actual, note)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: CreateContributionPlan :exec
INSERT INTO contribution_plan (uid, security_id, weight, amount, note)
VALUES (?, ?, ?, ?, ?);
