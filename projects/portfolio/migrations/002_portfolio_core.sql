-- +goose Up
CREATE TABLE IF NOT EXISTS category (
    id TEXT PRIMARY KEY DEFAULT (
        'cat_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (
        role IN (
            'liquidity',
            'growth',
            'income',
            'protection',
            'other'
        )
    ),
    parent_id TEXT DEFAULT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    UNIQUE (uid, name),
    FOREIGN KEY (parent_id) REFERENCES category (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS security (
    id TEXT PRIMARY KEY DEFAULT (
        'sec_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    symbol TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (
        type IN (
            'stock',
            'etf',
            'fund',
            'bond',
            'crypto',
            'other'
        )
    ),
    category_id TEXT DEFAULT NULL,
    currency TEXT NOT NULL DEFAULT 'KRW',
    note TEXT NOT NULL DEFAULT '',
    UNIQUE (uid, name),
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS account (
    id TEXT PRIMARY KEY DEFAULT (
        'acc_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    category_id TEXT NOT NULL,
    name TEXT NOT NULL,
    provider TEXT NOT NULL DEFAULT '',
    balance INTEGER NOT NULL DEFAULT 0,
    monthly_contrib INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS holding (
    id TEXT PRIMARY KEY DEFAULT (
        'h_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    security_id TEXT NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    target_amount INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    UNIQUE (uid, security_id),
    FOREIGN KEY (security_id) REFERENCES security (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS allocation_target (
    id TEXT PRIMARY KEY DEFAULT (
        'at_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    category_id TEXT NOT NULL,
    target_weight REAL NOT NULL DEFAULT 0 CHECK (
        target_weight >= 0
        AND target_weight <= 100
    ),
    target_amount INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    UNIQUE (uid, category_id),
    FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS budget_entry (
    id TEXT PRIMARY KEY DEFAULT (
        'be_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    name TEXT NOT NULL,
    direction TEXT NOT NULL CHECK (
        direction IN ('income', 'expense')
    ),
    class TEXT NOT NULL CHECK (
        class IN ('fixed', 'variable')
    ),
    planned INTEGER NOT NULL DEFAULT 0,
    actual INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    created TEXT DEFAULT CURRENT_TIMESTAMP,
    updated TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contribution_plan (
    id TEXT PRIMARY KEY DEFAULT (
        'cp_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    security_id TEXT NOT NULL,
    weight REAL NOT NULL DEFAULT 0 CHECK (
        weight >= 0
        AND weight <= 100
    ),
    amount INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    UNIQUE (uid, security_id),
    FOREIGN KEY (security_id) REFERENCES security (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_category_uid ON category (uid, display_order);

CREATE INDEX IF NOT EXISTS idx_account_uid ON account (uid);

CREATE INDEX IF NOT EXISTS idx_holding_uid ON holding (uid);

CREATE INDEX IF NOT EXISTS idx_budget_uid ON budget_entry (uid);

CREATE INDEX IF NOT EXISTS idx_plan_uid ON contribution_plan (uid);

CREATE VIEW IF NOT EXISTS v_total_asset AS
SELECT u.uid, COALESCE(
        (
            SELECT SUM(balance)
            FROM account a
            WHERE
                a.uid = u.uid
        ), 0
    ) + COALESCE(
        (
            SELECT SUM(h.amount)
            FROM holding h
            WHERE
                h.uid = u.uid
        ), 0
    ) AS total_amount
FROM (
        SELECT DISTINCT
            uid
        FROM account
        UNION
        SELECT DISTINCT
            uid
        FROM holding
    ) u;

CREATE VIEW IF NOT EXISTS v_category_allocation AS
SELECT
    c.uid,
    c.id AS category_id,
    c.name AS category_name,
    (
        COALESCE(acc.sum_balance, 0) + COALESCE(h.sum_holding, 0)
    ) AS amount,
    ROUND(
        100.0 * (
            COALESCE(acc.sum_balance, 0) + COALESCE(h.sum_holding, 0)
        ) / NULLIF(t.total_amount, 0),
        2
    ) AS weight_pct
FROM
    category c
    LEFT JOIN (
        SELECT category_id, uid, SUM(balance) AS sum_balance
        FROM account
        GROUP BY
            uid,
            category_id
    ) acc ON acc.category_id = c.id
    AND acc.uid = c.uid
    LEFT JOIN (
        SELECT s.category_id, h.uid, SUM(h.amount) AS sum_holding
        FROM holding h
            JOIN security s ON s.id = h.security_id
        GROUP BY
            h.uid,
            s.category_id
    ) h ON h.category_id = c.id
    AND h.uid = c.uid
    LEFT JOIN v_total_asset t ON t.uid = c.uid;

CREATE VIEW IF NOT EXISTS v_rebalance_gap AS
SELECT
    a.uid,
    a.category_id,
    a.target_weight,
    a.target_amount,
    v.amount AS current_amount,
    ROUND(
        (a.target_weight / 100.0) * t.total_amount
    ) AS target_by_weight,
    COALESCE(
        NULLIF(a.target_amount, 0),
        ROUND(
            (a.target_weight / 100.0) * t.total_amount
        )
    ) AS target_final,
    COALESCE(
        NULLIF(a.target_amount, 0),
        ROUND(
            (a.target_weight / 100.0) * t.total_amount
        )
    ) - v.amount AS gap_amount
FROM
    allocation_target a
    JOIN v_category_allocation v ON v.category_id = a.category_id
    AND v.uid = a.uid
    JOIN v_total_asset t ON t.uid = a.uid;

-- 샘플 데이터
INSERT OR IGNORE INTO
    user (uid, name, email)
VALUES (
        'uid_demo_portfolio',
        '김포트',
        'portfolio-demo@example.com'
    );

INSERT INTO
    category (
        id,
        uid,
        name,
        role,
        display_order
    )
VALUES (
        'cat_cash',
        'uid_demo_portfolio',
        '현금',
        'liquidity',
        1
    ),
    (
        'cat_savings',
        'uid_demo_portfolio',
        '예금',
        'liquidity',
        2
    ),
    (
        'cat_domestic_stock',
        'uid_demo_portfolio',
        '국내 주식',
        'growth',
        3
    ),
    (
        'cat_overseas_stock',
        'uid_demo_portfolio',
        '해외 주식',
        'growth',
        4
    ),
    (
        'cat_domestic_bond',
        'uid_demo_portfolio',
        '국내 채권',
        'income',
        5
    ),
    (
        'cat_overseas_bond',
        'uid_demo_portfolio',
        '해외 채권',
        'income',
        6
    ),
    (
        'cat_retirement',
        'uid_demo_portfolio',
        '연금',
        'protection',
        7
    ),
    (
        'cat_alternative',
        'uid_demo_portfolio',
        '대체 투자',
        'other',
        8
    );

INSERT INTO
    security (
        id,
        uid,
        symbol,
        name,
        type,
        category_id,
        currency,
        note
    )
VALUES (
        'sec_cash_fund',
        'uid_demo_portfolio',
        'CASH001',
        '머니마켓 펀드',
        'fund',
        'cat_cash',
        'KRW',
        '단기 유동성 자산'
    ),
    (
        'sec_kospi200',
        'uid_demo_portfolio',
        '069500',
        'KODEX 200',
        'etf',
        'cat_domestic_stock',
        'KRW',
        '국내 대표 지수 추종'
    ),
    (
        'sec_sp500',
        'uid_demo_portfolio',
        '295820',
        'TIGER 미국S&P500',
        'etf',
        'cat_overseas_stock',
        'KRW',
        '원화 환헤지 없음'
    ),
    (
        'sec_usbond',
        'uid_demo_portfolio',
        '448310',
        'TIGER 미국채10년선물',
        'etf',
        'cat_overseas_bond',
        'KRW',
        '장기채권 분산'
    ),
    (
        'sec_nasdaq',
        'uid_demo_portfolio',
        '133690',
        'TIGER 나스닥100',
        'etf',
        'cat_overseas_stock',
        'KRW',
        '성장주 비중 확대'
    ),
    (
        'sec_dividend',
        'uid_demo_portfolio',
        '314250',
        'SMART 배당주',
        'etf',
        'cat_domestic_stock',
        'KRW',
        '국내 고배당 ETF'
    ),
    (
        'sec_reit',
        'uid_demo_portfolio',
        'A372000',
        'ESR 켄달스퀘어 리츠',
        'stock',
        'cat_alternative',
        'KRW',
        '리츠 편입'
    ),
    (
        'sec_retirement_bond',
        'uid_demo_portfolio',
        'RET001',
        '퇴직연금 채권형',
        'fund',
        'cat_retirement',
        'KRW',
        '퇴직연금 기본형 상품'
    );

INSERT INTO
    account (
        id,
        uid,
        category_id,
        name,
        provider,
        balance,
        monthly_contrib,
        note
    )
VALUES (
        'acc_checking',
        'uid_demo_portfolio',
        'cat_cash',
        '입출금 통장',
        '카카오뱅크',
        2500000,
        1500000,
        '생활비 계좌'
    ),
    (
        'acc_savings',
        'uid_demo_portfolio',
        'cat_savings',
        '파킹 예금',
        '토스뱅크',
        32000000,
        500000,
        '비상금 및 단기 목표'
    ),
    (
        'acc_retire',
        'uid_demo_portfolio',
        'cat_retirement',
        '퇴직연금(우리)',
        '우리은행',
        48700000,
        500000,
        '퇴직연금 DC형'
    ),
    (
        'acc_domestic_bond',
        'uid_demo_portfolio',
        'cat_domestic_bond',
        '국내 채권 CMA',
        'NH투자',
        7800000,
        200000,
        '단기채 중심'
    ),
    (
        'acc_overseas_bond',
        'uid_demo_portfolio',
        'cat_overseas_bond',
        '달러 표시 채권형',
        '신한은행',
        6200000,
        200000,
        '환 노출 상품'
    ),
    (
        'acc_alternative',
        'uid_demo_portfolio',
        'cat_alternative',
        '부동산 펀드',
        'KB증권',
        9500000,
        0,
        '만기 2년 남음'
    );

INSERT INTO
    holding (
        id,
        uid,
        security_id,
        amount,
        target_amount,
        note
    )
VALUES (
        'h_kospi200',
        'uid_demo_portfolio',
        'sec_kospi200',
        28500000,
        30000000,
        ''
    ),
    (
        'h_sp500',
        'uid_demo_portfolio',
        'sec_sp500',
        31200000,
        32000000,
        ''
    ),
    (
        'h_usbond',
        'uid_demo_portfolio',
        'sec_usbond',
        8400000,
        9000000,
        '리밸런싱 예정'
    ),
    (
        'h_nasdaq',
        'uid_demo_portfolio',
        'sec_nasdaq',
        22800000,
        24000000,
        ''
    ),
    (
        'h_dividend',
        'uid_demo_portfolio',
        'sec_dividend',
        15600000,
        15000000,
        '분기 배당 재투자'
    ),
    (
        'h_reit',
        'uid_demo_portfolio',
        'sec_reit',
        7800000,
        8000000,
        '배당 수익 재투자'
    ),
    (
        'h_retirement',
        'uid_demo_portfolio',
        'sec_retirement_bond',
        50300000,
        52000000,
        '퇴직연금 비중 유지'
    );

INSERT INTO
    allocation_target (
        id,
        uid,
        category_id,
        target_weight,
        target_amount,
        note
    )
VALUES (
        'at_cash',
        'uid_demo_portfolio',
        'cat_cash',
        5,
        5000000,
        '생활비 2개월 분'
    ),
    (
        'at_savings',
        'uid_demo_portfolio',
        'cat_savings',
        20,
        30000000,
        ''
    ),
    (
        'at_domestic_stock',
        'uid_demo_portfolio',
        'cat_domestic_stock',
        25,
        38000000,
        ''
    ),
    (
        'at_overseas_stock',
        'uid_demo_portfolio',
        'cat_overseas_stock',
        25,
        38000000,
        '미국 주식 중심'
    ),
    (
        'at_overseas_bond',
        'uid_demo_portfolio',
        'cat_overseas_bond',
        10,
        15000000,
        ''
    ),
    (
        'at_retirement',
        'uid_demo_portfolio',
        'cat_retirement',
        10,
        20000000,
        ''
    ),
    (
        'at_alternative',
        'uid_demo_portfolio',
        'cat_alternative',
        5,
        8000000,
        '리츠 및 대체투자'
    );

INSERT INTO
    budget_entry (
        id,
        uid,
        name,
        direction,
        class,
        planned,
        actual,
        note
    )
VALUES (
        'be_salary',
        'uid_demo_portfolio',
        '급여',
        'income',
        'fixed',
        5500000,
        5500000,
        ''
    ),
    (
        'be_bonus',
        'uid_demo_portfolio',
        '성과급',
        'income',
        'variable',
        500000,
        600000,
        '분기 실적에 따라 변동'
    ),
    (
        'be_side',
        'uid_demo_portfolio',
        '부업 수입',
        'income',
        'variable',
        300000,
        250000,
        ''
    ),
    (
        'be_rent',
        'uid_demo_portfolio',
        '월세',
        'expense',
        'fixed',
        1200000,
        1200000,
        ''
    ),
    (
        'be_living',
        'uid_demo_portfolio',
        '생활비',
        'expense',
        'variable',
        900000,
        950000,
        '식비/교통/통신 포함'
    ),
    (
        'be_invest',
        'uid_demo_portfolio',
        '투자 적립',
        'expense',
        'fixed',
        1500000,
        1500000,
        ''
    ),
    (
        'be_travel',
        'uid_demo_portfolio',
        '여가/여행',
        'expense',
        'variable',
        300000,
        280000,
        ''
    ),
    (
        'be_education',
        'uid_demo_portfolio',
        '자기계발',
        'expense',
        'variable',
        200000,
        150000,
        '온라인 강의 구독'
    );

INSERT INTO
    contribution_plan (
        id,
        uid,
        security_id,
        weight,
        amount,
        note
    )
VALUES (
        'cp_kospi200',
        'uid_demo_portfolio',
        'sec_kospi200',
        25,
        350000,
        ''
    ),
    (
        'cp_sp500',
        'uid_demo_portfolio',
        'sec_sp500',
        25,
        350000,
        '환율 상황에 따라 조정'
    ),
    (
        'cp_nasdaq',
        'uid_demo_portfolio',
        'sec_nasdaq',
        20,
        280000,
        ''
    ),
    (
        'cp_dividend',
        'uid_demo_portfolio',
        'sec_dividend',
        15,
        210000,
        ''
    ),
    (
        'cp_usbond',
        'uid_demo_portfolio',
        'sec_usbond',
        10,
        140000,
        ''
    ),
    (
        'cp_reit',
        'uid_demo_portfolio',
        'sec_reit',
        5,
        70000,
        ''
    );

INSERT INTO
    casbin_rule (
        p_type,
        v0,
        v1,
        v2,
        v3,
        v4,
        v5
    )
VALUES (
        'p',
        'user',
        '*',
        '*',
        '',
        '',
        ''
    );

INSERT INTO
    casbin_rule (
        p_type,
        v0,
        v1,
        v2,
        v3,
        v4,
        v5
    )
VALUES (
        'p',
        'admin',
        '*',
        '*',
        '',
        '',
        ''
    );

INSERT INTO
    casbin_rule (
        p_type,
        v0,
        v1,
        v2,
        v3,
        v4,
        v5
    )
VALUES (
        'g',
        'uid_demo_portfolio',
        'admin',
        '',
        '',
        '',
        ''
    );

INSERT INTO
    casbin_rule (
        p_type,
        v0,
        v1,
        v2,
        v3,
        v4,
        v5
    )
VALUES (
        'g',
        'uid_demo_portfolio',
        'user',
        '',
        '',
        '',
        ''
    );

-- +goose Down
DROP VIEW IF EXISTS v_rebalance_gap;

DROP VIEW IF EXISTS v_category_allocation;

DROP VIEW IF EXISTS v_total_asset;

DROP INDEX IF EXISTS idx_plan_uid;

DROP INDEX IF EXISTS idx_budget_uid;

DROP INDEX IF EXISTS idx_holding_uid;

DROP INDEX IF EXISTS idx_account_uid;

DROP INDEX IF EXISTS idx_category_uid;

DROP TABLE IF EXISTS contribution_plan;

DROP TABLE IF EXISTS budget_entry;

DROP TABLE IF EXISTS allocation_target;

DROP TABLE IF EXISTS holding;

DROP TABLE IF EXISTS account;

DROP TABLE IF EXISTS security;

DROP TABLE IF EXISTS category;