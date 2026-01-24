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
| 環境変数 | `NEXT_PUBLIC_API_BASE_URL` |

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

| 環境 | 設定ファイル | API URL |
|------|-------------|---------|
| 開発 | `.env.local` | `http://localhost:8080` |
| 本番 | `.env.production` | `https://api.sikoulab.com` |

```bash
# frontend/.env.local（開発環境）
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# frontend/.env.production（本番環境）
NEXT_PUBLIC_API_BASE_URL=https://api.sikoulab.com
```

### ファイル配置

```
frontend/
├── lib/
│   └── api-client.ts    # axiosクライアント
├── types/
│   └── api/             # API型定義
│       ├── article.ts
│       ├── user.ts
│       └── index.ts
└── .env.local           # 環境変数（git管理外）
```

## 実装例

### axiosクライアント

```typescript
// frontend/lib/api-client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
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
w.Header().Set("Access-Control-Allow-Origin", "https://sikoulab.com")
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
