-- plansテーブルから amount / currency カラムを削除する
ALTER TABLE plans
    DROP COLUMN IF EXISTS amount,
    DROP COLUMN IF EXISTS currency;
