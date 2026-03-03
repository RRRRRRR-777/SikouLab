-- 'inactive' ユーザーを 'trialing' に戻す
UPDATE users
SET subscription_status = 'trialing', updated_at = NOW()
WHERE subscription_status = 'inactive';

-- デフォルト値を削除
ALTER TABLE users ALTER COLUMN subscription_status DROP DEFAULT;

-- 制約を元に戻す
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_subscription_status;
ALTER TABLE users ADD CONSTRAINT check_subscription_status
    CHECK (subscription_status IN ('active', 'canceled', 'past_due', 'trialing'));
