import axios from "axios";

/**
 * APIクライアントのaxiosインスタンス。
 * 環境変数NEXT_PUBLIC_API_URLからベースURLを取得する。
 */
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

// レスポンスエラーのインターセプター
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // APIエラーをコンソールに記録（本番ではエラートラッキングに送信）
    console.error("API Error:", error.response?.data ?? error.message);
    return Promise.reject(error);
  }
);
