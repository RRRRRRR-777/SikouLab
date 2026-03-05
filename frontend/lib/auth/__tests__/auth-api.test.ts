/**
 * 認証APIクライアントのテスト
 *
 * @description
 * Firebase認証とバックエンドAPIの連携をテストする。
 * - ログイン（ID Token送信）
 * - セッション確認
 * - ログアウト
 *
 * @see {@link file://./../../../docs/functions/auth/login.md} 認証機能仕様
 */

import { describe, it, expect, vi, beforeEach } from "vitest";

// vi.hoisted()でモックインスタンスを先にホイスティング
const { mockAxiosInstance } = vi.hoisted(() => {
  return {
    mockAxiosInstance: {
      post: vi.fn(),
      get: vi.fn(),
    },
  };
});

// 共有apiClientをモック
vi.mock("@/lib/api", () => ({
  apiClient: mockAxiosInstance,
}));

import { authApi } from "../auth-api";

describe("authApi", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("login", () => {
    it("成功時はユーザー情報と初回ログインフラグを返す", async () => {
      // 準備: モックレスポンス
      const mockResponse = {
        data: {
          user: {
            id: 1,
            oauth_provider: "google.com",
            name: "Test User",
            display_name: "Test User",
            avatar_url: "https://example.com/avatar.jpg",
            role: "user",
            plan_id: null,
            subscription_status: null,
            created_at: "2024-01-01T00:00:00Z",
            updated_at: "2024-01-01T00:00:00Z",
          },
          is_first_login: false,
        },
      };
      mockAxiosInstance.post.mockResolvedValue(mockResponse);

      // 実行
      const result = await authApi.login("valid_id_token");

      // 検証
      expect(result.user).toEqual(mockResponse.data.user);
      expect(result.is_first_login).toBe(false);
      expect(mockAxiosInstance.post).toHaveBeenCalledWith("/auth/login", {
        id_token: "valid_id_token",
      });
    });

    it("初回ログイン時はis_first_loginがtrue", async () => {
      // 準備: 初回ログインのモックレスポンス
      const mockResponse = {
        data: {
          user: {
            id: 1,
            oauth_provider: "google.com",
            name: "New User",
            display_name: "New User",
            avatar_url: null,
            role: "user",
            plan_id: null,
            subscription_status: null,
            created_at: "2024-01-01T00:00:00Z",
            updated_at: "2024-01-01T00:00:00Z",
          },
          is_first_login: true,
        },
      };
      mockAxiosInstance.post.mockResolvedValue(mockResponse);

      // 実行
      const result = await authApi.login("new_user_token");

      // 検証
      expect(result.is_first_login).toBe(true);
    });

    it("無効なID Token時は401エラーを投げる", async () => {
      // 準備: 認証エラーのモック
      const error = {
        response: {
          status: 401,
          data: {
            code: "UNAUTHORIZED",
            message: "無効なID Tokenです",
          },
        },
      };
      mockAxiosInstance.post.mockRejectedValue(error);

      // 実行・検証
      await expect(authApi.login("invalid_token")).rejects.toMatchObject({
        response: { status: 401 },
      });
    });

    it("ネットワークエラー時はエラーを投げる", async () => {
      // 準備: ネットワークエラーのモック
      const error = new Error("Network Error");
      mockAxiosInstance.post.mockRejectedValue(error);

      // 実行・検証
      await expect(authApi.login("any_token")).rejects.toThrow("Network Error");
    });
  });

  describe("getMe", () => {
    it("成功時はユーザー情報を返す", async () => {
      // 準備: モックレスポンス
      const mockResponse = {
        data: {
          user: {
            id: 1,
            oauth_provider: "google.com",
            name: "Test User",
            display_name: "Test User",
            avatar_url: "https://example.com/avatar.jpg",
            role: "user",
            plan_id: 1,
            subscription_status: "active",
            created_at: "2024-01-01T00:00:00Z",
            updated_at: "2024-01-01T00:00:00Z",
          },
        },
      };
      mockAxiosInstance.get.mockResolvedValue(mockResponse);

      // 実行
      const result = await authApi.getMe();

      // 検証
      expect(result.user).toEqual(mockResponse.data.user);
      expect(mockAxiosInstance.get).toHaveBeenCalledWith("/auth/me");
    });

    it("未認証時は401エラーを投げる", async () => {
      // 準備: 未認証エラーのモック
      const error = {
        response: {
          status: 401,
          data: {
            code: "UNAUTHORIZED",
            message: "未認証です",
          },
        },
      };
      mockAxiosInstance.get.mockRejectedValue(error);

      // 実行・検証
      await expect(authApi.getMe()).rejects.toMatchObject({
        response: { status: 401 },
      });
    });
  });

  describe("logout", () => {
    it("成功時は204 No Contentを返す", async () => {
      // 準備: モックレスポンス
      const mockResponse = {
        status: 204,
        data: null,
      };
      mockAxiosInstance.post.mockResolvedValue(mockResponse);

      // 実行
      await authApi.logout();

      // 検証
      expect(mockAxiosInstance.post).toHaveBeenCalledWith("/auth/logout");
    });

    it("未認証時は401エラーを投げる", async () => {
      // 準備: 未認証エラーのモック
      const error = {
        response: {
          status: 401,
          data: {
            code: "UNAUTHORIZED",
            message: "未認証です",
          },
        },
      };
      mockAxiosInstance.post.mockRejectedValue(error);

      // 実行・検証
      await expect(authApi.logout()).rejects.toMatchObject({
        response: { status: 401 },
      });
    });
  });
});
