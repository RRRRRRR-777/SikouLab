package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RRRRRRR-777/SicouLab/backend/internal/domain"
)

// contextKey はcontextキーの衝突を避けるための非公開型。
// 文字列型をそのまま使うと他パッケージのキーと衝突する可能性があるため、独自型を使用する。
type contextKey string

// userContextKey はcontextにUserを格納するためのキー。
const userContextKey contextKey = "user"

// roleHierarchy はロールの階層を定義する。
// 上位ロールほど高い値を持ち、includesRole の比較に使用する。
var roleHierarchy = map[string]int{
	"user":   1,
	"writer": 2,
	"admin":  3,
}

// authSessionVerifier はセッション検証インターフェース。
// AuthUsecaseへの直接依存を避け、テスト時のモック差し替えを可能にする。
type authSessionVerifier interface {
	GetCurrentUser(ctx context.Context, sessionToken string) (*domain.User, error)
}

// authErrorResponse はミドルウェアが返すエラーレスポンスのJSON構造。
type authErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeAuthError はミドルウェア用のJSONエラーレスポンスを書き込む。
func writeAuthError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(authErrorResponse{Code: code, Message: message})
}

// RequireAuth は認証済みユーザーのみにアクセスを制限するためのミドルウェアを返す。
//
// Cookieなし・トークン検証失敗・ユーザー不在の場合は全て401を返す。
// 検証成功時はcontextにUserを設定してnextハンドラーを呼ぶ。
func RequireAuth(uc authSessionVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// セッションCookieが存在しない場合は認証不可
			cookie, err := r.Cookie("session")
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
				return
			}

			// セッショントークンを検証してユーザーを取得
			user, err := uc.GetCurrentUser(r.Context(), cookie.Value)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
				return
			}

			// トークンは有効でもDBにユーザーが存在しない場合は認証不可
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
				return
			}

			// contextにユーザーを設定して次のハンドラーへ
			ctx := ContextWithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// hasRequiredRole はuserRoleがallowedRolesのいずれかを満たすか判定する。
//
// ロール階層（admin > writer > user）を考慮し、
// 上位ロールは下位ロールのアクセス権限を包含する。
func hasRequiredRole(userRole string, allowedRoles []string) bool {
	userLevel, ok := roleHierarchy[userRole]
	if !ok {
		// 未定義ロールはアクセス不可
		return false
	}

	for _, allowed := range allowedRoles {
		allowedLevel, ok := roleHierarchy[allowed]
		if !ok {
			continue
		}
		// userLevelがallowedLevel以上であればアクセス可（上位ロールは下位権限を包含）
		if userLevel >= allowedLevel {
			return true
		}
	}
	return false
}

// RequireRole は指定されたroleのみ通すミドルウェア。
//
// ロール階層（admin > writer > user）を考慮する。
// allowedRolesに"writer"を指定した場合、adminも通過できる。
// contextにUserがない場合は401、ロール不足は403を返す。
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// RequireAuth後に呼ばれることを前提とするが、Userがない場合は401
			user := UserFromContext(r.Context())
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
				return
			}

			// ロール階層に基づいてアクセス権を判定
			if !hasRequiredRole(user.Role, roles) {
				writeAuthError(w, http.StatusForbidden, "FORBIDDEN", "この操作を行う権限がありません")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSubscription はサブスク有効ユーザーのみ通すミドルウェア。
//
// subscription_status が "active" のユーザーのみ通過できる。
// contextにUserがない場合は401、サブスク無効は403を返す。
func RequireSubscription() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// RequireAuth後に呼ばれることを前提とするが、Userがない場合は401
			user := UserFromContext(r.Context())
			if user == nil {
				writeAuthError(w, http.StatusUnauthorized, "UNAUTHORIZED", "認証が必要です")
				return
			}

			// "active" 以外のサブスクリプション状態はアクセス不可
			if user.SubscriptionStatus != "active" {
				writeAuthError(w, http.StatusForbidden, "FORBIDDEN", "有効なサブスクリプションが必要です")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserFromContext はハンドラーが認証済みユーザー情報を参照するために使用する。
//
// RequireAuthを通過したリクエストのcontextに格納されたユーザーを返す。
// contextにUserが設定されていない場合はnilを返す。
func UserFromContext(ctx context.Context) *domain.User {
	u, _ := ctx.Value(userContextKey).(*domain.User)
	return u
}

// ContextWithUser はテストや内部処理でユーザー情報を直接注入するために使用する。
//
// 通常のリクエストフローではRequireAuthが自動的に設定する。
func ContextWithUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
