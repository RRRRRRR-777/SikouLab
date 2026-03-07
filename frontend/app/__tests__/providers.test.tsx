/**
 * Providersコンポーネントのテスト
 *
 * @description
 * Providersコンポーネントに AuthProvider が含まれていること、
 * および TanStack Query の QueryClientProvider が含まれていることを確認する。
 */

import { describe, it, expect, vi } from "vitest";
import { renderHook, waitFor } from "@testing-library/react";
import { Providers } from "../providers";
import { useAuth } from "../../lib/auth/auth-context";
import { useQueryClient } from "@tanstack/react-query";

// vi.hoisted() でモックを先にホイスティング
const { mockUnsubscribe } = vi.hoisted(() => {
  return {
    mockUnsubscribe: vi.fn(),
  };
});

// Next.jsルーターのモック
vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
  }),
  usePathname: () => "/",
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

// Firebase認証モジュールのモック
// AuthProviderが Providers 内に存在する場合、onAuthStateChangedHelper が呼ばれ isLoading が false になる
vi.mock("../../lib/auth/firebase", () => ({
  signInWithGoogle: vi.fn(),
  signInWithApple: vi.fn(),
  signInWithX: vi.fn(),
  signOut: vi.fn(),
  onAuthStateChangedHelper: vi.fn((callback: (user: null) => void) => {
    // 初期状態は未認証（null）を渡す
    callback(null);
    return mockUnsubscribe;
  }),
  getIdToken: vi.fn(),
}));

// 認証APIのモック
vi.mock("../../lib/auth/auth-api", () => ({
  authApi: {
    login: vi.fn(),
    getMe: vi.fn(),
    logout: vi.fn(),
  },
}));

describe("Providers", () => {
  it("AuthProviderが含まれており、useAuthが使える", async () => {
    /**
     * Providers コンポーネント内で useAuth を呼び出したとき、
     * AuthProvider が内包されていれば onAuthStateChangedHelper が呼ばれ isLoading が false になる。
     * @param root0 - コンポーネントのプロパティ
     * @param root0.children - 子コンポーネント
     * @returns Providersでラップされた子要素
     */
    const wrapper = ({ children }: { children: React.ReactNode }) => (
      <Providers>{children}</Providers>
    );

    const { result } = renderHook(() => useAuth(), { wrapper });

    // AuthProvider が含まれていれば、onAuthStateChangedHelper が null を渡してくるため
    // isLoading は false になる（含まれていなければデフォルト値 true のまま）
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });
  });

  it("QueryClientProviderが含まれている", () => {
    /**
     * Providers コンポーネント内で TanStack Query の useQueryClient が利用可能なことを確認する。
     * QueryClientProvider は既に providers.tsx に含まれているため、このテストは通過する。
     * @param root0 - コンポーネントのプロパティ
     * @param root0.children - 子コンポーネント
     * @returns Providersでラップされた子要素
     */
    const wrapper = ({ children }: { children: React.ReactNode }) => (
      <Providers>{children}</Providers>
    );

    const { result } = renderHook(() => useQueryClient(), { wrapper });

    // QueryClientProvider内に配置されていれば useQueryClient が値を返す
    expect(result.current).toBeDefined();
  });
});
