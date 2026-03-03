/**
 * 認証コンテキストのテスト
 *
 * @description
 * 認証状態管理と認証操作をテストする。
 * - 認証状態の初期化
 * - ログイン処理
 * - ログアウト処理
 * - ユーザー情報の取得
 *
 * @see {@link file://./../../../docs/functions/auth/login.md} 認証機能仕様
 */

import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import type { User } from "../types";

// vi.hoisted()でモックインスタンスを先にホイスティング
const { mockPush, mockAuthApi, mockFirebaseUser, mockUnsubscribe, mockPathnameRef } = vi.hoisted(() => {
  const firebaseUser = {
    uid: "firebase_uid_123",
    email: "test@example.com",
    displayName: "Test User",
    photoURL: "https://example.com/avatar.jpg",
    getIdToken: vi.fn(() => Promise.resolve("valid_firebase_id_token")),
  };

  return {
    mockPush: vi.fn(),
    mockPathnameRef: { current: "/" },
    mockAuthApi: {
      login: vi.fn(() => Promise.resolve({
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
      })),
      getMe: vi.fn(() => Promise.resolve({
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
      })),
      logout: vi.fn(() => Promise.resolve()),
    },
    mockFirebaseUser: firebaseUser,
    mockUnsubscribe: vi.fn(),
  };
});

// Next.jsルーターのモック
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
    replace: vi.fn(),
  }),
  usePathname: () => mockPathnameRef.current,
}));

// sonnerのモック
vi.mock("sonner", () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Firebase Authのモック
vi.mock("firebase/auth", async () => {
  const firebase = await vi.importActual<typeof import("firebase/auth")>("firebase/auth");

  return {
    ...firebase,
    onAuthStateChanged: vi.fn(() => mockUnsubscribe),
  };
});

// Firebaseユーザーのモック型
type MockFirebaseUser = {
  uid: string;
  email: string | null;
  displayName: string | null;
  photoURL: string | null;
};

// Firebase認証モジュールのモック
vi.mock("../firebase", () => ({
  signInWithGoogle: vi.fn(() => Promise.resolve(mockFirebaseUser)),
  signInWithApple: vi.fn(() => Promise.resolve(mockFirebaseUser)),
  signInWithX: vi.fn(() => Promise.resolve(mockFirebaseUser)),
  signOut: vi.fn(() => Promise.resolve()),
  onAuthStateChangedHelper: vi.fn(
    (callback: (user: MockFirebaseUser | null) => void) => {
      // 初期状態は未認証（null）を渡す
      callback(null);
      return mockUnsubscribe;
    }
  ),
  getIdToken: vi.fn(() =>
    Promise.resolve("valid_firebase_id_token")
  ),
}));

// 認証APIのモック
vi.mock("../auth-api", () => ({
  authApi: mockAuthApi,
}));

import { useAuth } from "../auth-context";
import { AuthProvider } from "../auth-context";
import { onAuthStateChangedHelper, getIdToken } from "../firebase";

describe("useAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPathnameRef.current = "/";
  });

  describe("初期状態", () => {
    it("認証状態は未ログイン・ローディング完了", () => {
      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      expect(result.current.user).toBeNull();
      expect(result.current.isLoading).toBe(false);
      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe("loginWithGoogle", () => {
    it("成功時はユーザー情報をセットする", async () => {
      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.loginWithGoogle();
      });

      // 検証
      await waitFor(() => {
        expect(result.current.user).toEqual({
          id: 1,
          oauthProvider: "google.com",
          name: "Test User",
          displayName: "Test User",
          avatarUrl: "https://example.com/avatar.jpg",
          role: "user",
          planId: null,
          subscriptionStatus: null,
          createdAt: "2024-01-01T00:00:00Z",
          updatedAt: "2024-01-01T00:00:00Z",
        });
        expect(result.current.isAuthenticated).toBe(true);
        expect(mockAuthApi.login).toHaveBeenCalledWith("valid_firebase_id_token");
      });
    });

    it("初回ログイン時はサブスクリプション画面へ遷移する", async () => {
      mockAuthApi.login.mockResolvedValueOnce({
        user: {
          id: 2,
          oauth_provider: "google.com",
          name: "New User",
          display_name: "New User",
          avatar_url: null,
          role: "user",
          plan_id: null,
          subscription_status: null,
          created_at: "2024-01-01T00:00:00Z",
          updated_at: "2024-01-01T00:00:00Z",
        } satisfies User,
        is_first_login: true,
      } as never);

      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.loginWithGoogle();
      });

      // 検証: サブスクリプション画面へ遷移
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith("/subscription");
      });
    });
  });

  describe("loginWithApple", () => {
    it("成功時はユーザー情報をセットする", async () => {
      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.loginWithApple();
      });

      // 検証
      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
        expect(mockAuthApi.login).toHaveBeenCalled();
      });
    });
  });

  describe("loginWithX", () => {
    it("成功時はユーザー情報をセットする", async () => {
      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.loginWithX();
      });

      // 検証
      await waitFor(() => {
        expect(result.current.isAuthenticated).toBe(true);
        expect(mockAuthApi.login).toHaveBeenCalled();
      });
    });
  });

  describe("logout", () => {
    it("成功時はユーザー情報をクリアする", async () => {
      mockAuthApi.logout.mockResolvedValue(undefined);

      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.logout();
      });

      // 検証
      await waitFor(() => {
        expect(result.current.user).toBeNull();
        expect(result.current.isAuthenticated).toBe(false);
        expect(mockAuthApi.logout).toHaveBeenCalled();
      });
    });

  });

  describe("ログインページでの動作", () => {
    it("Firebaseユーザーがいても/auth/meを呼ばない", async () => {
      // ログインページのパスを設定
      mockPathnameRef.current = "/login";

      // FirebaseユーザーありでonAuthStateChangedを呼ぶ
      vi.mocked(onAuthStateChangedHelper).mockImplementationOnce((callback) => {
        callback(mockFirebaseUser as never);
        return mockUnsubscribe;
      });

      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });

      // ログインページでは/auth/meも/auth/loginも呼ばれないことを確認
      expect(mockAuthApi.getMe).not.toHaveBeenCalled();
      expect(mockAuthApi.login).not.toHaveBeenCalled();
      expect(result.current.user).toBeNull();
    });
  });

  describe("refresh", () => {
    it("成功時は最新のユーザー情報を取得する", async () => {
      const { result } = renderHook(() => useAuth(), {
        wrapper: AuthProvider,
      });

      // 実行
      await act(async () => {
        await result.current.refresh();
      });

      // 検証
      await waitFor(() => {
        expect(mockAuthApi.getMe).toHaveBeenCalled();
      });
    });
  });

  describe("onAuthStateChanged リフレッシュ", () => {
    beforeEach(() => {
      // Firebaseユーザーあり・ログインページ以外でonAuthStateChangedを呼ぶ
      vi.mocked(onAuthStateChangedHelper).mockImplementationOnce((callback) => {
        callback(mockFirebaseUser as never);
        return mockUnsubscribe;
      });
    });

    it("Firebaseユーザーがいる場合、getIdTokenが呼ばれてからPOST /auth/loginが呼ばれる", async () => {
      renderHook(() => useAuth(), { wrapper: AuthProvider });

      // 検証: getIdToken(forceRefresh=true) → authApi.login の順で呼ばれる
      await waitFor(() => {
        expect(vi.mocked(getIdToken)).toHaveBeenCalledWith(mockFirebaseUser, true);
        expect(mockAuthApi.login).toHaveBeenCalledWith("valid_firebase_id_token");
      });
    });

    it("トークンリフレッシュ成功後、loginレスポンスからユーザー情報がセットされる", async () => {
      const { result } = renderHook(() => useAuth(), { wrapper: AuthProvider });

      // 検証: loginレスポンスのユーザー情報がセットされる
      await waitFor(() => {
        expect(result.current.user).toEqual({
          id: 1,
          oauthProvider: "google.com",
          name: "Test User",
          displayName: "Test User",
          avatarUrl: "https://example.com/avatar.jpg",
          role: "user",
          planId: null,
          subscriptionStatus: null,
          createdAt: "2024-01-01T00:00:00Z",
          updatedAt: "2024-01-01T00:00:00Z",
        });
        expect(result.current.isAuthenticated).toBe(true);
        expect(result.current.isLoading).toBe(false);
      });
    });

    it("getIdToken失敗時はユーザーをnullにする", async () => {
      // 準備: getIdTokenが失敗するようにモック
      vi.mocked(getIdToken).mockRejectedValueOnce(new Error("Token refresh failed"));

      const { result } = renderHook(() => useAuth(), { wrapper: AuthProvider });

      // 検証: ユーザーはnull、ローディング完了
      await waitFor(() => {
        expect(result.current.user).toBeNull();
        expect(result.current.isAuthenticated).toBe(false);
        expect(result.current.isLoading).toBe(false);
      });
    });

    it("login失敗（401等）時はユーザーをnullにする", async () => {
      // 準備: loginが401エラーを返すようにモック
      mockAuthApi.login.mockRejectedValueOnce({ response: { status: 401 } });

      const { result } = renderHook(() => useAuth(), { wrapper: AuthProvider });

      // 検証: ユーザーはnull、ローディング完了
      await waitFor(() => {
        expect(result.current.user).toBeNull();
        expect(result.current.isAuthenticated).toBe(false);
        expect(result.current.isLoading).toBe(false);
      });
    });
  });
});
