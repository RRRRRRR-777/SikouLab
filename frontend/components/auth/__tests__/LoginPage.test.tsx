/**
 * ログイン画面のテスト
 *
 * @description
 * ログイン画面（S-01）のUIと操作をテストする。
 * - 画面の表示
 * - OAuthプロバイダボタンのクリック
 * - アクセシビリティ
 *
 * @see {@link file://./../../../docs/functions/auth/login.md} 認証機能仕様
 * @see {@link file://./../../../docs/versions/1_0_0/SicouLab.pen} Pencilデザイン（h3Lxa）
 */

import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { LoginPage } from "../LoginPage";

// 認証フックのモック
const mockLoginWithGoogle = vi.fn();
const mockLoginWithApple = vi.fn();
const mockLoginWithX = vi.fn();

vi.mock("../../../lib/auth/auth-context", () => ({
  useAuth: () => ({
    loginWithGoogle: mockLoginWithGoogle,
    loginWithApple: mockLoginWithApple,
    loginWithX: mockLoginWithX,
    isLoading: false,
    isAuthenticated: false,
  }),
}));

// Next.js Linkのモック
vi.mock("next/link", () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}));

// Next.js Imageのモック
vi.mock("next/image", () => ({
  // eslint-disable-next-line @next/next/no-img-element -- テストコード内のモック
  default: ({ src, alt }: { src: string; alt: string }) => <img src={src} alt={alt} />,
}));

describe("LoginPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("画面表示", () => {
    it("ロゴとタイトルが表示される", () => {
      render(<LoginPage />);

      // ロゴテキスト
      expect(screen.getByText("シコウラボ")).toBeInTheDocument();
      // 利用規約文言
      expect(
        screen.getByText(/ログインすることで利用規約に同意したものとみなされます/),
      ).toBeInTheDocument();
    });

    it("3つのOAuthボタンが表示される", () => {
      render(<LoginPage />);

      expect(screen.getByRole("button", { name: /Google でログイン/ })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /Apple でログイン/ })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /X でログイン/ })).toBeInTheDocument();
    });
  });

  describe("OAuthボタン操作", () => {
    it("Googleボタンクリック時はGoogleログインを実行する", () => {
      render(<LoginPage />);

      const googleButton = screen.getByRole("button", { name: /Google でログイン/ });
      fireEvent.click(googleButton);

      expect(mockLoginWithGoogle).toHaveBeenCalledTimes(1);
    });

    it("Appleボタンクリック時はAppleログインを実行する", () => {
      render(<LoginPage />);

      const appleButton = screen.getByRole("button", { name: /Apple でログイン/ });
      fireEvent.click(appleButton);

      expect(mockLoginWithApple).toHaveBeenCalledTimes(1);
    });

    it("Xボタンクリック時はXログインを実行する", () => {
      render(<LoginPage />);

      const xButton = screen.getByRole("button", { name: /X でログイン/ });
      fireEvent.click(xButton);

      expect(mockLoginWithX).toHaveBeenCalledTimes(1);
    });
  });

  describe("アクセシビリティ", () => {
    it("各OAuthボタンには適切なtype属性が設定されている", () => {
      render(<LoginPage />);

      const googleButton = screen.getByRole("button", { name: /Google でログイン/ });
      const appleButton = screen.getByRole("button", { name: /Apple でログイン/ });
      const xButton = screen.getByRole("button", { name: /X でログイン/ });

      expect(googleButton).toHaveAttribute("type", "button");
      expect(appleButton).toHaveAttribute("type", "button");
      expect(xButton).toHaveAttribute("type", "button");
    });
  });

  describe("ダークモード対応", () => {
    it("ダークモードでも正しく表示される", () => {
      // ダークモードを設定
      document.documentElement.classList.add("dark");

      render(<LoginPage />);

      // 基本要素が表示されることを確認
      expect(screen.getByText("シコウラボ")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /Google でログイン/ })).toBeInTheDocument();

      // 後処理
      document.documentElement.classList.remove("dark");
    });
  });
});
