-- plansテーブルのシードデータ
-- 注意: このファイルは開発環境でのみ使用すること

-- ベースプラン（既存データはINSERT ... ON CONFLICTで重複回避）
INSERT INTO plans (name, description, is_active)
VALUES ('ベースプラン', '記事・ニュース・銘柄ページの閲覧が可能なプラン', true)
ON CONFLICT DO NOTHING;

-- 将来の上位プラン（コメントアウト済み。必要時に有効化）
-- INSERT INTO plans (name, description, is_active)
-- VALUES ('プレミアムプラン', '全ての機能が利用可能なプラン', true)
-- ON CONFLICT DO NOTHING;
