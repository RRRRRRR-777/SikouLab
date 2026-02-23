-- usersテーブル: ユーザー情報・認証・課金状態
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    oauth_provider VARCHAR(50) NOT NULL,
    oauth_user_id VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    display_name VARCHAR(100),
    avatar_url TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    plan_id BIGINT REFERENCES plans(id) ON DELETE SET NULL,
    univapay_customer_id VARCHAR(255) UNIQUE,
    subscription_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- プロバイダ内での一意性制約
    UNIQUE(oauth_provider, oauth_user_id)
);

-- roleカラムのチェック制約
ALTER TABLE users ADD CONSTRAINT check_role
    CHECK (role IN ('admin', 'writer', 'user'));

-- subscription_statusカラムのチェック制約
ALTER TABLE users ADD CONSTRAINT check_subscription_status
    CHECK (subscription_status IN ('active', 'canceled', 'past_due', 'trialing'));

-- インデックス作成
CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_user_id);
CREATE INDEX idx_users_plan_id ON users(plan_id);
CREATE INDEX idx_users_subscription_status ON users(subscription_status);
