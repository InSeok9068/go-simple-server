-- +goose Up
CREATE TABLE IF NOT EXISTS category (
    id TEXT PRIMARY KEY DEFAULT (
        'cat_' || LOWER(HEX(RANDOMBLOB(7)))
    ),
    uid TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (
        role IN (
            'liquidity', -- 유동성 자산 (예 : 현금, 예금, CMA)
            'growth', -- 성장 자산 (예 : 주식, 주식형 ETF, 부동산, 벤처 투자)
            'income', -- 소득 자산 (예 : 채권, 배당주, 리츠(REITs), 채권형 ETF)
            'protection', -- 보호 자산 (예 : 보험상품, 금, 달러, 채권)
            'other' -- 기타 자산 (예 : 암호화폐, 미술품, 수집품 등)
        )
    ),
    parent_id TEXT DEFAULT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    UNIQUE (uid, name),
    FOREIGN KEY (parent_id) REFERENCES category (id) ON DELETE SET NULL
);