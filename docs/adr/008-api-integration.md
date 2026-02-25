# ADR-008: API連携方針

## ステータス

採用

## 背景

フロントエンド（Next.js）とバックエンド（Go）の通信方式を定義する必要がある。セキュリティを最優先としつつ、将来のモバイルアプリ対応も考慮する。

## 決定事項

### API通信

| 項目 | 決定 |
|------|------|
| ライブラリ | axios |
| タイムアウト | 10秒 |
| リトライ | TanStack Queryに任せる |
| 環境変数 | `NEXT_PUBLIC_API_URL` |

**理由:**
- axiosのインターセプターで認証・エラー処理を一元化
- TanStack Queryのデフォルトリトライ（3回）で十分

### 認証方式

#### Web（現在）

| 項目 | 設定 |
|------|------|
| 方式 | JWT + httpOnly Cookie |
| HttpOnly | true（XSS対策） |
| Secure | true（HTTPS必須） |
| SameSite | Lax（CSRF対策） |
| MaxAge | 7日間 |

**理由:**
- セキュリティ最優先
- SameSite=LaxでCSRF対策が簡潔に実現

#### Mobile（将来対応時）

| 項目 | 設定 |
|------|------|
| 方式 | JWT + Authorization header |
| 保存先 | セキュアストレージ |
| Go API | Cookie/Header両方式をサポート |

### 型定義

| 項目 | 決定 |
|------|------|
| 方式 | OpenAPIから手動でTypeScript型を作成 |
| 命名規則 | camelCase統一（Go側も） |
| 管理場所 | `frontend/types/api/` |

**理由:**
- 自動生成ツールを避け、実装をシンプルに保つ
- 命名変換処理が不要

### エラーハンドリング

共通ラッパーで一元化:
- 401 → ログイン画面へリダイレクト
- 500 → エラートースト表示

### 環境変数

| 環境 | 設定ファイル | NEXT_PUBLIC_API_URL |
|------|-------------|---------------------|
| 開発 | `.env`（git管理外） | `http://localhost:8080/api/v1` |
| 本番 | 未設定 | -（LBルーティングを使用） |

```bash
# frontend/.env（開発環境・git管理外）
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# 本番環境は未設定。Cloud Load Balancing の URLパスルーティングに委ねる。
# ブラウザから /api/v1/* → LB → Cloud Run（Go API）
```

### 本番環境のBFFアーキテクチャ

本番環境では Next.js が BFF（Backend for Frontend）として機能する。Cloud Load Balancing が URL パスでルーティングを担当する（インフラ設計書 §5.2 参照）:

| パス | ルーティング先 |
|------|-------------|
| `/api/*` | Cloud Run（Go API） |
| `/*` | Cloud Run（Next.js） |

同一ドメインでフロントエンドとバックエンドを提供するため、ブラウザから見て同一オリジンとなり CORS 設定が不要になる。

> 開発環境は直接接続（クロスオリジン）のため、バックエンドの CORS 設定（`Access-Control-Allow-Credentials: true`）が必要。

### ファイル配置

```
frontend/
├── lib/
│   ├── api.ts           # 汎用axiosクライアント
│   └── auth/
│       └── auth-api.ts  # 認証API専用クライアント
├── types/
│   └── api/             # API型定義
│       ├── article.ts
│       ├── user.ts
│       └── index.ts
└── .env                 # 環境変数（git管理外）
```

## 実装例

### axiosクライアント

```typescript
// frontend/lib/api.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "/api/v1",
  timeout: 10000,
  withCredentials: true, // Cookie送信
});

// リクエストインターセプター
apiClient.interceptors.request.use((config) => {
  config.headers['Content-Type'] = 'application/json';
  return config;
});

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

### Go API（将来のモバイル対応）

```go
// 認証ミドルウェア
func GetToken(r *http.Request) string {
    // Authorization headerを優先（モバイル）
    if authHeader := r.Header.Get("Authorization"); authHeader != "" {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }

    // Cookieをチェック（Web）
    if cookie, err := r.Cookie("token"); err == nil {
        return cookie.Value
    }

    return ""
}
```

### Cookie設定（Go）

```go
http.SetCookie(w, &http.Cookie{
    Name:     "token",
    Value:    jwtToken,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   60 * 60 * 24 * 7, // 7日間
    Path:     "/",
})
```

### CORS設定（Go）

```go
w.Header().Set("Access-Control-Allow-Origin", "https://sicoulab.com")
w.Header().Set("Access-Control-Allow-Credentials", "true")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

## 影響範囲

- フロントエンド: 全API呼び出し
- バックエンド: 認証ミドルウェア、CORS設定
- インフラ: 環境変数設定

## 関連

- [ADR-003: データフェッチ](./003-data-fetching.md)
- [開発ガイドライン](../development_guidelines.md)
