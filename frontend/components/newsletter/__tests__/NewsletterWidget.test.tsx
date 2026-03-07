/**
 * メール登録ウィジェットのテスト
 *
 * @description
 * 機能要件5（メールアドレス登録）のウィジェットコンポーネント単体テスト。
 * 表示・バリデーション・API呼び出し・トースト表示をテストする。
 *
 * @see {@link file://../../../../docs/functions/settings/home.md} 詳細設計書
 */

import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { NewsletterWidget } from "../NewsletterWidget";

// --- モック変数のホイスティング ---
const {
  mockSubscribeNewsletter,
  mockGetNewsletterSubscription,
  mockToastSuccess,
  mockToastError,
  mockUseAuth,
} = vi.hoisted(() => ({
  mockSubscribeNewsletter: vi.fn(),
  mockGetNewsletterSubscription: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
  mockUseAuth: vi.fn(),
}));

// --- 依存モジュールのモック ---

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

vi.mock("@/lib/settings/settings-api", () => ({
  settingsApi: {
    subscribeNewsletter: (...args: unknown[]) =>
      mockSubscribeNewsletter(...args),
    getNewsletterSubscription: (...args: unknown[]) =>
      mockGetNewsletterSubscription(...args),
  },
}));

vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => mockUseAuth(),
}));

// --- テストデータ ---

const mockSubscription = {
  id: 1,
  email: "test@example.com",
  is_active: true,
  created_at: "2026-03-04T00:00:00Z",
  updated_at: "2026-03-04T00:00:00Z",
};

// --- テスト ---

describe("NewsletterWidget", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // デフォルト: ログイン済み・購読未登録
    mockUseAuth.mockReturnValue({
      user: { id: 1, name: "テストユーザー" },
      isAuthenticated: true,
      isLoading: false,
    });
    mockGetNewsletterSubscription.mockRejectedValue({ response: { status: 404 } });
  });

  // ===== 表示テスト =====

  describe("ウィジェット表示", () => {
    it("ログイン済みの場合、タイトル・説明文・入力欄・ボタンが表示される", async () => {
      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(screen.getByText("メールで記事を受け取る")).toBeInTheDocument();
      });
      expect(
        screen.getByText(/毎朝7:30に記事の要約をお届けします/),
      ).toBeInTheDocument();
      expect(
        screen.getByPlaceholderText("メールアドレス"),
      ).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: "購読する" }),
      ).toBeInTheDocument();
    });

    it("未ログインの場合、ウィジェットが表示されない", () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });

      const { container } = render(<NewsletterWidget />);

      expect(container.innerHTML).toBe("");
    });

    it("認証ローディング中はウィジェットが表示されない", () => {
      mockUseAuth.mockReturnValue({
        user: null,
        isAuthenticated: false,
        isLoading: true,
      });

      const { container } = render(<NewsletterWidget />);

      expect(container.innerHTML).toBe("");
    });

    it("既に購読済みの場合は「購読済みです」が表示される", async () => {
      mockGetNewsletterSubscription.mockResolvedValue(mockSubscription);

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(screen.getByText("購読済みです")).toBeInTheDocument();
      });
      // 入力欄とボタンは表示されない
      expect(
        screen.queryByPlaceholderText("メールアドレス"),
      ).not.toBeInTheDocument();
      expect(
        screen.queryByRole("button", { name: "購読する" }),
      ).not.toBeInTheDocument();
    });
  });

  // ===== API呼び出しテスト =====

  describe("購読登録", () => {
    it("メールアドレスを入力して購読ボタンをクリックするとAPIが呼ばれる", async () => {
      mockSubscribeNewsletter.mockResolvedValue(mockSubscription);

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockSubscribeNewsletter).toHaveBeenCalledWith({ email: "user@example.com" });
      });
    });

    it("成功時にトーストが表示される", async () => {
      mockSubscribeNewsletter.mockResolvedValue(mockSubscription);

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockToastSuccess).toHaveBeenCalledWith("購読を開始しました");
      });
    });

    it("成功後に「購読済みです」表示に切り替わる", async () => {
      mockSubscribeNewsletter.mockResolvedValue(mockSubscription);

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(screen.getByText("購読済みです")).toBeInTheDocument();
      });
    });
  });

  // ===== バリデーションテスト =====

  describe("バリデーション", () => {
    it("空文字で購読ボタンをクリックするとバリデーションエラー", async () => {
      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "購読する" }),
        ).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockSubscribeNewsletter).not.toHaveBeenCalled();
      });
      expect(
        screen.getByText("メールアドレスを入力してください"),
      ).toBeInTheDocument();
    });

    it("不正なメール形式でバリデーションエラー", async () => {
      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "invalid-email" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockSubscribeNewsletter).not.toHaveBeenCalled();
      });
      expect(
        screen.getByText("有効なメールアドレスを入力してください"),
      ).toBeInTheDocument();
    });
  });

  // ===== エラーハンドリングテスト =====

  describe("エラーハンドリング", () => {
    it("APIエラー時にエラートーストが表示される", async () => {
      mockSubscribeNewsletter.mockRejectedValue({
        response: {
          status: 400,
          data: { message: "不正なメールアドレスです" },
        },
      });

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockToastError).toHaveBeenCalledWith(
          "不正なメールアドレスです",
        );
      });
    });

    it("APIエラー（メッセージなし）時にデフォルトエラーメッセージが表示される", async () => {
      mockSubscribeNewsletter.mockRejectedValue(new Error("Network Error"));

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(mockToastError).toHaveBeenCalledWith("購読登録に失敗しました");
      });
    });

    it("送信中はボタンが非活性になる", async () => {
      // subscribe がpendingで止まる
      mockSubscribeNewsletter.mockReturnValue(new Promise(() => {}));

      render(<NewsletterWidget />);

      await waitFor(() => {
        expect(
          screen.getByPlaceholderText("メールアドレス"),
        ).toBeInTheDocument();
      });

      fireEvent.change(screen.getByPlaceholderText("メールアドレス"), {
        target: { value: "user@example.com" },
      });
      fireEvent.click(screen.getByRole("button", { name: "購読する" }));

      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "購読する" }),
        ).toBeDisabled();
      });
    });
  });
});
