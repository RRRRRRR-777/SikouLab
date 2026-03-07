-- usersテーブルにemailカラムを追加
-- OAuthプロバイダから取得したメールアドレスを保存する（読み取り専用、設定画面で表示）
ALTER TABLE users ADD COLUMN email VARCHAR(255);
