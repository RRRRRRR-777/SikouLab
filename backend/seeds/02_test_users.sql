-- テスト用ユーザーのシードデータ
-- 注意: このファイルは開発環境でのみ使用すること

-- テスト用管理者ユーザー（Google OAuth想定）
INSERT INTO users (
    oauth_provider,
    oauth_user_id,
    name,
    display_name,
    avatar_url,
    role,
    plan_id,
    subscription_status
)
VALUES (
    'google',
    'test_admin_001',
    'Test Admin',
    'テスト管理者',
    'https://example.com/avatar/admin.png',
    'admin',
    1, -- ベースプラン
    'active'
)
ON CONFLICT (oauth_provider, oauth_user_id) DO NOTHING;

-- テスト用ライターユーザー
INSERT INTO users (
    oauth_provider,
    oauth_user_id,
    name,
    display_name,
    avatar_url,
    role,
    plan_id,
    subscription_status
)
VALUES (
    'google',
    'test_writer_001',
    'Test Writer',
    'テストライター',
    'https://example.com/avatar/writer.png',
    'writer',
    1, -- ベースプラン
    'active'
)
ON CONFLICT (oauth_provider, oauth_user_id) DO NOTHING;

-- テスト用一般ユーザー
INSERT INTO users (
    oauth_provider,
    oauth_user_id,
    name,
    display_name,
    avatar_url,
    role,
    plan_id,
    subscription_status
)
VALUES (
    'google',
    'test_user_001',
    'Test User',
    'テストユーザー',
    'https://example.com/avatar/user.png',
    'user',
    1, -- ベースプラン
    'active'
)
ON CONFLICT (oauth_provider, oauth_user_id) DO NOTHING;

-- テスト用未課金ユーザー（初回ログイン後、サブスク未登録の状態）
INSERT INTO users (
    oauth_provider,
    oauth_user_id,
    name,
    display_name,
    avatar_url,
    role,
    subscription_status
)
VALUES (
    'google',
    'test_inactive_001',
    'Test Inactive User',
    'テスト未課金ユーザー',
    'https://example.com/avatar/inactive.png',
    'user',
    'inactive'
)
ON CONFLICT (oauth_provider, oauth_user_id) DO NOTHING;

-- テスト用ユーザー設定
INSERT INTO user_settings (user_id, sidebar_article_expanded, sidebar_admin_expanded)
SELECT id, true, false FROM users WHERE oauth_user_id LIKE 'test_%'
ON CONFLICT (user_id) DO NOTHING;
