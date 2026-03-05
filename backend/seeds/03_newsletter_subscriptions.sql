-- newsletter_subscriptionsテーブルのシードデータ
-- 注意: このファイルは開発環境でのみ使用すること

-- テスト管理者のニュースレター購読（有効）
INSERT INTO newsletter_subscriptions (user_id, email, is_active)
SELECT id, 'admin@example.com', true FROM users WHERE oauth_user_id = 'test_admin_001'
ON CONFLICT (user_id) DO NOTHING;

-- テストライターのニュースレター購読（有効）
INSERT INTO newsletter_subscriptions (user_id, email, is_active)
SELECT id, 'writer@example.com', true FROM users WHERE oauth_user_id = 'test_writer_001'
ON CONFLICT (user_id) DO NOTHING;

-- テスト一般ユーザーのニュースレター購読（無効 = 購読解除済み）
INSERT INTO newsletter_subscriptions (user_id, email, is_active)
SELECT id, 'user@example.com', false FROM users WHERE oauth_user_id = 'test_user_001'
ON CONFLICT (user_id) DO NOTHING;
