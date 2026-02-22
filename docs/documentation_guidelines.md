# ドキュメント戦略

## 1. 目的
本ドキュメントは、SicouLabプロジェクトにおけるドキュメント作成・管理・運用の方針を定義する。
開発者およびAIエージェントが迷わず、統一された方法でドキュメントを扱えるようにすることを目的とする。

## 2. ドキュメント構造

### 2.1 ディレクトリ構成
```
docs/
├── SUMMARY.md                                 # ドキュメントインデックス（必須）
├── _sidebar.md                                # Docsifyサイドバー定義
├── development_guidelines.md                  # 開発ガイドライン
├── documentation_guidelines.md                # 本ドキュメント
├── versions/                                  # バージョン別ドキュメント
│   └── 1_0_0/                                # バージョン1.0.0
│       ├── requirements.md                    # 要件定義書
│       ├── system-design.md                   # 基本設計書
│       ├── data-model.md                      # データ設計書（ER図・テーブル定義）
│       ├── infrastructure.md                  # インフラ設計書
│       └── development-plan.md                # 開発計画書
├── functions/                                 # 機能別詳細仕様
│   ├── auth/                                  # 認証機能
│   ├── article/                               # 記事機能
│   ├── news/                                  # ニュース機能
│   ├── stock/                                 # 銘柄機能
│   ├── poll/                                  # アンケート機能
│   └── ...
├── adr/                                       # Architecture Decision Records
│   └── NNN-title.md                          # ADRドキュメント
└── daily/                                     # 議事録
    └── MMDD.md                               # 日付ベースの議事録
```

### 2.2 各ドキュメントの役割

#### SUMMARY.md（インデックス）
- **目的**: docs配下の全ドキュメントの概要を一覧で把握
- **対象読者**: 開発者・AIエージェント
- **更新タイミング**: docs配下のマークダウンファイル追加・更新時（daily/配下を除く）
- **配置場所**: docs直下

#### development_guidelines.md（開発ガイドライン）
- **目的**: 技術スタック、CI/CD、テスト方針などの開発ルールを定義
- **対象読者**: 開発者
- **更新タイミング**: 開発プロセス変更時
- **配置場所**: docs直下

#### documentation_guidelines.md（本ドキュメント）
- **目的**: ドキュメント作成・管理・運用の方針を定義
- **対象読者**: 開発者
- **更新タイミング**: ドキュメント戦略変更時
- **配置場所**: docs直下

#### _sidebar.md（サイドバー）
- **目的**: Docsifyのナビゲーションサイドバーを定義
- **対象読者**: ドキュメントサイト閲覧者
- **更新タイミング**: docs配下のドキュメント追加・リネーム・削除時
- **配置場所**: docs直下

#### versions/（バージョン別ドキュメント）
- **目的**: バージョンごとの要件定義・設計書を管理
- **命名規則**: `a_b_c`形式（例: 1_0_0）
- **対象読者**: 全ステークホルダー
- **含まれるドキュメント**:
  - requirements.md（要件定義書）
  - system-design.md（基本設計書）
  - data-model.md（データ設計書）
  - infrastructure.md（インフラ設計書）
  - development-plan.md（開発計画書）

#### functions/（機能別詳細仕様）
- **目的**: 機能ごとの詳細仕様を管理
- **粒度**: 中機能単位（例: 記事作成、記事詳細など）
- **対象読者**: 開発者
- **更新タイミング**: 機能実装時・仕様変更時
- **ディレクトリ構成**: 大機能（auth, article等）配下に中機能ごとのファイルを配置

#### adr/（Architecture Decision Records）
- **目的**: 重要なアーキテクチャ上の意思決定を記録
- **命名規則**: `NNN-title.md`（例: 001-use-clean-architecture.md）
- **対象読者**: 開発者、アーキテクト
- **更新タイミング**: アーキテクチャ上の重要な決定時

#### daily/（議事録）
- **目的**: 日々の議論や意思決定の記録
- **命名規則**: `MMDD.md`（例: 0107.md）
- **対象読者**: プロジェクトメンバー
- **更新タイミング**: ミーティング後
- **注意**: SUMMARY.mdへの記載は不要

## 3. ファイル命名規則

### 3.1 基本ルール
- **形式**: ケバブケース（小文字、ハイフン区切り）
- **言語**: 英語
- **拡張子**: `.md`（マークダウン形式）

### 3.2 特殊なケース
- **議事録**: `MMDD.md`（月日の数字4桁）
- **ADR**: `NNN-title.md`（連番3桁 + ハイフン + タイトル）
- **バージョンディレクトリ**: `X_Y_Z`（アンダースコア区切り、例: 1_0_0）
- **機能ディレクトリ**: 大機能名（小文字）配下に中機能ファイル（ケバブケース）

### 3.3 例
```
○ docs/functions/auth/login.md
○ docs/functions/article/create-edit.md
× docs/functions/認証/ログイン.md
× docs/functions/auth/createEdit.md
```

## 4. ドキュメント更新ルール

### 4.1 SUMMARY.md更新ルール（必須）
docs配下のマークダウンファイル（*.md）が追加・更新された際は、**必ず**SUMMARY.mdを更新する。
（daily/配下の議事録は除く）

**記載形式:**
```markdown
* ファイルパス
  * 概要: 2〜3行でそのファイルの内容が簡潔に正確にわかる説明
```

### 4.2 _sidebar.md更新ルール（必須）
docs配下にドキュメントを**追加・リネーム・削除**した際は、**必ず**_sidebar.mdも更新する。
- ファイルパスは実際のファイル名と正確に一致させること
- 「機能別詳細仕様」セクションは機能番号（F-XX）の昇順で並べること

### 4.3 情報探索ルール
1. まずSUMMARY.mdの概要を確認
2. 概要から情報に辿り着きそうであれば、該当するファイルを参照
3. AIエージェントは常にこのフローに従う

### 4.4 ドキュメント作成時の注意
- **目的を明確に**: 何のためのドキュメントかを冒頭に記載
- **対象読者を意識**: 誰が読むかを想定して記述
- **最新性を保つ**: 実装と乖離しないよう定期的に更新
- **TBD（To Be Determined）を活用**: 未決定事項は明示的にTBDと記載

## 5. 設計ドキュメントワークフロー

機能の設計・開発は以下の順序で進める。

### 5.1 ワークフロー

```
/requirements
  └─ docs/versions/<version>/requirements.md
       └─ 機能一覧に機能を定義
            ↓
/system-design
  └─ docs/versions/<version>/system-design.md
            ↓
/feature-spec
  └─ docs/functions/<feature-category>/<feature-name>.md
       └─ system-design.md の機能一覧のステータス・リンクを更新
            ↓
実装
```

### 5.2 対応Skill

| 作業 | 呼び出すSkill |
|------|-------------|
| 要件定義書の新規作成・大幅修正 | `/requirements` |
| 基本設計書の新規作成・大幅修正 | `/system-design` |
| 詳細設計書の新規作成・大幅修正 | `/feature-spec` |
| 作業見積もりの追加 | `/feature-spec-estimate-cost` |

### 5.3 運用ルール
- 詳細設計を作成したら、基本設計の機能一覧を更新する（ステータス・リンク）
- 新規機能は原則として基本設計に追記してから詳細設計を行う
- 緊急対応など例外的に詳細設計を先行する場合は、後から基本設計に反映する

## 6. ER図同期ルール

### 6.1 単一正（Single Source of Truth）
テーブル定義（カラム名・データ型・PK/FK・説明）は `docs/versions/<version>/data-model.md` で一元管理する。

### 6.2 詳細設計書のER図ルール
- 詳細設計書のER図は、data-model.mdから**関連テーブルの全カラム定義をコピー**する
- リレーションは機能に関連するもののみ記載する
- ER図の冒頭に正への参照コメントを記載:
  ```
  %% 正: docs/versions/1_0_0/data-model.md
  ```

### 6.3 同期フロー
- **data-model.md → function docs**: data-model.mdを変更した場合、影響を受ける全ての詳細設計書の関連テーブルセクションを更新する
- **function docs → data-model.md**: 詳細設計で新規テーブル/カラムが必要になった場合、**先にdata-model.mdを更新**してから詳細設計に反映する

### 6.4 禁止事項
- 詳細設計書でdata-model.mdに存在しないカラムを独自定義すること
- 詳細設計書でカラム名・データ型をdata-model.mdと異なる形で記載すること

## 7. ドキュメント間の相関関係

| ドキュメント | ファイルパス | 概要 |
|---|---|---|
| 要件定義書 | `docs/versions/<version>/requirements.md` | Why/What。背景・目的・成功基準・機能一覧 |
| 基本設計書 | `docs/versions/<version>/system-design.md` | How（外部設計）。機能相関・画面設計・権限 |
| データ設計書 | `docs/versions/<version>/data-model.md` | ER図・テーブル定義の単一正 |
| 詳細設計書 | `docs/functions/<category>/<name>.md` | How（内部設計）。機能単位の詳細仕様・シーケンス図 |

### 編集時の影響範囲

| 編集ファイル | 確認・修正対象 |
|---|---|
| requirements.md | system-design.md |
| system-design.md | requirements.md、functions/配下の全詳細設計 |
| data-model.md | functions/配下の関連する全詳細設計 |
| functions/*.md | system-design.md（機能一覧のステータス・リンク） |

## 8. バージョン管理
- **Git管理**: 全ドキュメントはGitで管理
- **コミットメッセージ**: ドキュメント変更内容を明記（日本語）
- **レビュー**: TBD（今後プルリクエストベースのレビューフローを検討）

## 9. AIエージェント向けガイドライン
- 情報探索は`docs/SUMMARY.md`を起点にすること
- ドキュメント追加・更新時は必ず`SUMMARY.md`と`_sidebar.md`を更新すること
- 概要は簡潔かつ正確に記述すること
- TBD項目は明示的に記載し、曖昧な表現を避けること
- 詳細ルールは`.claude/rules/docs.md`を参照すること

## 10. 今後の検討事項
- **ドキュメントレビュープロセス**: TBD
- **ドキュメントテンプレート整備**: TBD
- **ドキュメント自動生成**: TBD（OpenAPIからの自動生成など）
