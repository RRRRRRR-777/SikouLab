-- newsletter_subscriptionsテーブル: ニュースレター購読管理
CREATE TABLE IF NOT EXISTS newsletter_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- インデックス作成（user_idはUNIQUE制約で自動インデックスが作成されるため不要）
CREATE INDEX idx_newsletter_subscriptions_is_active ON newsletter_subscriptions(is_active);
