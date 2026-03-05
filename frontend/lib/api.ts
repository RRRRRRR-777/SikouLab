import axios from "axios";

/**
 * APIクライアントのaxiosインスタンス。
 *
 * 開発環境では NEXT_PUBLIC_API_URL（例: http://localhost:8080/api/v1）を設定する。
 * 本番環境では未設定とし、Cloud Load Balancing の URL パスルーティングに委ねる。
 *
 * withCredentials: true によりCookieが同一オリジンで送信される（セッション認証に必要）。
 */
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

// レスポンスエラーのインターセプター
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // 認証エラー（401）は呼び出し元で処理するためログ出力をスキップ
    if (error.response?.status !== 401) {
      console.error("API Error:", error.response?.data ?? error.message);
    }
    return Promise.reject(error);
  },
);
