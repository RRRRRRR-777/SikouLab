---
paths:
  - "docs/**/*.md"
---

# docsディレクトリルール

## CLAUDE.md更新ルール
* **追加**: 繰り返し説明する内容、禁止コマンド、プロダクト原則
* **追加しない**: 技術詳細はrulesへ
* **簡潔に**: 指示が増えると精度が下がる

---

## SUMMARY.md更新ルール
* docs配下(docs/\*)のマークダウンファイル(\*.md)が更新された際は`docs/SUMMARY.md`を**必ず更新**する
* docs/SUMMARY.mdの記載ルール:
  * 概要はそのファイルの内容が簡潔に正確にわかるものを記載
    ```
    * ファイルパス
      * 概要: 2、3行
    ```
  * docs/daily配下のファイルは議事録であるので不要です

## _sidebar.md更新ルール
* docs配下に新規ドキュメントを追加・リネーム・削除した際は`docs/_sidebar.md`を**必ず更新**する
* `.claude/`配下の`.md`ファイルを追加・リネーム・削除した際は`docs/_sidebar.md`を**必ず更新**する
* ファイルパスは実際のファイル名と正確に一致させること

## 情報参照の優先順位
* わからないことがあれば`docs/SUMMARY.md`を最初に確認すること
  * 概要を読んで情報に辿り着きそうであればそのファイルパスのファイルを参照する

## ドキュメント作成方針
* 読みやすさを最優先
* 階層は深くしすぎない（最大3階層推奨）
* 技術的詳細よりも「なぜそうするか」を重視

## ファイル・ディレクトリ命名規則
* **ファイル名・ディレクトリ名は英語**（ケバブケース推奨）
* 内容（本文）は日本語でOK
* 例:
  * ○ `docs/functions/auth/login.md`
  * ○ `docs/functions/article/create-edit.md`
  * × `docs/functions/認証/ログイン.md`

## デザイン・UI参照
画面の設計を行う際は、以下のPencilファイルを参照すること。

* Pencil: `docs/versions/1_0_0/SikouLab.pen`

**注意事項**
* Pencilはあくまでモックアップであり、デザインが完全に正しいとは限らない
* 一般的なUI/UXのベストプラクティスも意識すること
* Pencilのデザインとベストプラクティスのどちらを採用すべきか判断に迷う場合は、ユーザーに確認すること

## 設計ドキュメントワークフロー

機能の設計・開発は以下の順序で進める。

### 1. 要件定義（/requirements）
* **Skill**: `/requirements`
* **出力先**: `docs/versions/<version>/requirements.md`

### 2. 基本設計（/system-design）
* **Skill**: `/system-design`
* **出力先**: `docs/versions/<version>/system-design.md`

### 3. 詳細設計（/feature-spec）
* **Skill**: `/feature-spec`
* **出力先**: `docs/functions/<feature-category>/<feature-name>.md`（英語ケバブケース）
  * 例: `docs/functions/auth/login.md`, `docs/functions/article/create-edit.md`

### ワークフロー図

```
/requirements
  └─ docs/versions/<version>/requirements.md
       └─ 機能一覧に機能を定義
            ↓
/feature-spec
  └─ docs/functions/<feature-category>/<feature-name>.md
       └─ 機能一覧のステータス・リンクを更新
            ↓
実装
```

### 運用ルール
* 詳細設計を作成したら、基本設計の機能一覧を更新する（ステータス・リンク）
* 新規機能は原則として基本設計に追記してから詳細設計を行う
* 緊急対応など例外的に詳細設計を先行する場合は、後から基本設計に反映する

## ドキュメント間の相関関係マッピング

### ドキュメント名とファイルパスの対応

| ドキュメント名 | ファイルパス | 概要 |
|---|---|---|
| **要件定義書** | `docs/versions/<version>/requirements.md` | Why/What を定義。背景・目的・成功基準・機能一覧 |
| **基本設計書** | `docs/versions/<version>/system-design.md` | How（外部設計）を定義。機能相関・画面設計・データ設計・権限 |
| **E2Eシナリオ** | `docs/versions/<version>/test_scenarios.md` | ユーザーストーリー・画面遷移に基づくE2Eテストシナリオ |
| **詳細設計書** | `docs/functions/<feature-category>/<feature-name>.md` | How（内部設計）を定義。機能単位の詳細仕様・シーケンス図 |

### ファイル相関図

```
要件定義書 (requirements.md)
    ↓↑ 常に同期
基本設計書 (system-design.md)
    │
    ├─→ E2Eシナリオ (test_scenarios.md)
    │
    ↓↑ フィードバック
詳細設計書 (docs/functions/<feature-category>/<feature-name>.md)
```

### 編集時の影響範囲

#### 1. `requirements.md` を編集した場合

**確認・修正対象**: `docs/versions/1_0_0/system-design.md`

影響を受ける箇所：
- 2.1 機能階層 - 機能の追加/削除/優先度変更
- 2.2 機能相関図 - 機能依存関係の変更
- 2.3 機能依存関係 - 機能間の依存内容の変更
- 3. 画面設計 - 画面一覧・画面遷移
- 4. データ設計 - テーブル定義・ER図
- 6. 権限マトリクス - ロール・権限の変更
- 7. 制約・前提 - ビジネス制約・プラン構成

#### 2. `system-design.md` を編集した場合

**確認・修正対象**:
1. `docs/versions/1_0_0/requirements.md`
2. `docs/functions/` 配下の**全ての詳細設計ファイル**

**requirements.md の影響を受ける箇所:**
- 5. 機能一覧 - 機能の追加/削除/優先度/ステータス/詳細設計リンク

**functions/ 配下の影響を受ける箇所:**
- **機能階層の変更** → 対応する詳細設計ファイルの追加/削除/統合
- **機能依存関係の変更** → 影響を受ける詳細設計の関連機能セクションを更新
- **画面設計の変更** → 対応する詳細設計ファイルの画面設計図セクションを更新
- **データ設計の変更** → 対応する詳細設計ファイルの関連テーブルセクションを更新
- **権限マトリクスの変更** → 対応する詳細設計ファイルの権限セクションを更新

#### 3. `docs/functions/<feature-category>/<feature-name>.md` を編集した場合

**確認・修正対象**: `docs/versions/1_0_0/system-design.md`

影響を受ける箇所：
- 2.1 機能階層 - ステータス・詳細設計リンクの更新
- 2.2 機能相関図 - 機能依存関係の修正（影響範囲の拡大検出時）
- 2.3 機能依存関係 - 依存内容の修正（新しい依存を発見した場合）
- 3. 画面設計 - 画面遷移の修正（新規画面追加時）
- 4. データ設計 - テーブル定義の追加/修正（新規テーブル発見時）

### 相関関係マトリクス

| 編集ファイル | requirements.md | system-design.md | functions/ |
|---|:---:|:---:|:---:|
| requirements.md | - | ○ 必ず修正 | ○ 間接的に確認 |
| system-design.md | ○ 必ず修正 | - | ○ 必ず修正 |
| functions/ | ✓ フィードバック | ✓ フィードバック | - |

凡例：
- ○: 確認・修正が必要
- ✓: 矛盾がないかレビュー推奨
- -: 該当なし
