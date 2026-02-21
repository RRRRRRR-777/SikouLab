-- plansテーブル: サブスクリプションプランのマスタデータ
CREATE TABLE IF NOT EXISTS plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 初期プランデータ（v1.0.0はベースプランのみ。上位プランは将来バージョンで追加）
INSERT INTO plans (name, description, is_active) VALUES
    ('ベースプラン', '記事・ニュース・銘柄ページの閲覧が可能なプラン', true);
