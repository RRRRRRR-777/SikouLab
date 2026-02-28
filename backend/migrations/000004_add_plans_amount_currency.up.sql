-- plansテーブルに amount / currency カラムを追加する
-- amount: 月額料金（最小通貨単位。例: JPY の場合は円）
-- currency: 通貨コード（ISO-4217。例: JPY）
ALTER TABLE plans
    ADD COLUMN amount INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN currency VARCHAR(3) NOT NULL DEFAULT 'JPY';
