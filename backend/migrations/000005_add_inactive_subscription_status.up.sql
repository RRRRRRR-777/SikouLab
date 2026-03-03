-- subscription_statusに'inactive'を追加し、デフォルト値を設定する。
-- 既存の CHECK 制約を DROP して再作成する。

-- 既存の制約を削除
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_subscription_status;

-- 'inactive' を含む新しい制約を追加
ALTER TABLE users ADD CONSTRAINT check_subscription_status
    CHECK (subscription_status IN ('active', 'canceled', 'past_due', 'trialing', 'inactive'));

-- デフォルト値を 'inactive' に設定
ALTER TABLE users ALTER COLUMN subscription_status SET DEFAULT 'inactive';

-- 既存の未課金ユーザー（univapay_customer_id が NULL かつ active）を inactive に更新
UPDATE users
SET subscription_status = 'inactive', updated_at = NOW()
WHERE univapay_customer_id IS NULL AND subscription_status = 'active';
