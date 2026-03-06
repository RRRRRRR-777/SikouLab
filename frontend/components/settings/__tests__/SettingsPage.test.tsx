/**
 * 設定画面（S-17）の統合テスト
 *
 * @description
 * 設定画面の4セクション（プロフィール、サブスクリプション、メール、FAQ）の
 * 表示と操作をテストする。
 *
 * @see {@link file://../../../../docs/functions/settings/home.md} 詳細設計書
 */

import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { SettingsPage } from "../SettingsPage";

// --- モック変数のホイスティング ---
const {
  mockToastSuccess,
  mockToastError,
  mockUpdateProfile,
  mockUploadAvatar,
  mockDeleteAvatar,
  mockGetSubscription,
  mockCreatePortal,
  mockGetNewsletterSubscription,
  mockSubscribeNewsletter,
  mockUnsubscribeNewsletter,
  mockUpdateNewsletterEmail,
  mockRefresh,
} = vi.hoisted(() => ({
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
  mockUpdateProfile: vi.fn(),
  mockUploadAvatar: vi.fn(),
  mockDeleteAvatar: vi.fn(),
  mockGetSubscription: vi.fn(),
  mockCreatePortal: vi.fn(),
  mockGetNewsletterSubscription: vi.fn(),
  mockSubscribeNewsletter: vi.fn(),
  mockUnsubscribeNewsletter: vi.fn(),
  mockUpdateNewsletterEmail: vi.fn(),
  mockRefresh: vi.fn(),
}));

// --- 依存モジュールのモック ---

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => "/settings",
}));

vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({
    user: {
      id: 1,
      oauthProvider: "google.com",
      oauthUserId: "firebase-uid-123",
      name: "yamada_taro",
      displayName: "山田 太郎",
      avatarUrl: "https://example.com/avatar.jpg",
      role: "user",
      planId: 1,
      subscriptionStatus: "active",
      createdAt: "2026-01-01T00:00:00Z",
      updatedAt: "2026-01-01T00:00:00Z",
    },
    isAuthenticated: true,
    isLoading: false,
    refresh: mockRefresh,
  }),
}));

vi.mock("@/lib/settings/settings-api", () => ({
  settingsApi: {
    updateProfile: (...args: unknown[]) => mockUpdateProfile(...args),
    uploadAvatar: (...args: unknown[]) => mockUploadAvatar(...args),
    deleteAvatar: (...args: unknown[]) => mockDeleteAvatar(...args),
    getSubscription: (...args: unknown[]) => mockGetSubscription(...args),
    createPortal: (...args: unknown[]) => mockCreatePortal(...args),
    getNewsletterSubscription: (...args: unknown[]) => mockGetNewsletterSubscription(...args),
    subscribeNewsletter: (...args: unknown[]) => mockSubscribeNewsletter(...args),
    unsubscribeNewsletter: (...args: unknown[]) => mockUnsubscribeNewsletter(...args),
    updateNewsletterEmail: (...args: unknown[]) => mockUpdateNewsletterEmail(...args),
  },
}));

// Next.js Imageのモック
vi.mock("next/image", () => ({
  default: ({ src, alt, ...props }: { src: string; alt: string; [key: string]: unknown }) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={src} alt={alt} {...props} />
  ),
}));

// --- テストデータ ---

const mockSubscriptionData = {
  plan_name: "プレミアムプラン",
  amount: 1980,
  currency: "JPY",
  status: "active" as const,
};

const mockNewsletterData = {
  id: 1,
  email: "newsletter@example.com",
  is_active: true,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

// --- テスト ---

describe("SettingsPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // デフォルトの正常系レスポンスをセットアップ
    mockGetSubscription.mockResolvedValue(mockSubscriptionData);
    mockGetNewsletterSubscription.mockResolvedValue(mockNewsletterData);
  });

  // ===== 設定画面全体の表示 =====

  describe("設定画面全体の表示", () => {
    it("4つのセクションが全て表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("プロフィール設定")).toBeInTheDocument();
        expect(screen.getByText("サブスクリプション")).toBeInTheDocument();
        expect(screen.getByText("メールアドレス設定")).toBeInTheDocument();
        expect(screen.getByText("FAQ・問い合わせ")).toBeInTheDocument();
      });
    });

    it("ページタイトル「設定」が表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByRole("heading", { name: "設定", level: 1 })).toBeInTheDocument();
      });
    });
  });

  // ===== F-10-1: プロフィールセクション =====

  describe("プロフィールセクション（F-10-1）", () => {
    it("表示名が表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByDisplayValue("山田 太郎")).toBeInTheDocument();
      });
    });

    it("ユーザーIDが読み取り専用で表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("firebase-uid-123")).toBeInTheDocument();
      });
    });

    it("アバター画像が表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        const avatar = screen.getByAltText("アバター");
        expect(avatar).toBeInTheDocument();
        expect(avatar).toHaveAttribute("src", "https://example.com/avatar.jpg");
      });
    });

    it("表示名を変更して保存できる", async () => {
      mockUpdateProfile.mockResolvedValue({
        user: {
          id: 1,
          name: "yamada_taro",
          display_name: "新しい表示名",
          avatar_url: "https://example.com/avatar.jpg",
          role: "user",
          plan_id: 1,
          subscription_status: "active",
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
        },
      });

      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByDisplayValue("山田 太郎")).toBeInTheDocument();
      });

      // 表示名入力フィールドを変更
      const input = screen.getByDisplayValue("山田 太郎");
      fireEvent.change(input, { target: { value: "新しい表示名" } });

      // プロフィールセクション内の保存ボタンをクリック
      const profileSection = screen.getByText("プロフィール設定").closest("section")!;
      const saveButtons = profileSection.querySelectorAll("button");
      const saveButton = Array.from(saveButtons).find((btn) => btn.textContent === "保存")!;
      fireEvent.click(saveButton);

      await waitFor(() => {
        expect(mockUpdateProfile).toHaveBeenCalledWith({
          display_name: "新しい表示名",
        });
        expect(mockToastSuccess).toHaveBeenCalledWith("保存しました");
      });
    });

    it("アバター変更ボタンが表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByRole("button", { name: /変更/ })).toBeInTheDocument();
      });
    });

    it("アバター削除ボタンで削除APIが呼ばれる", async () => {
      mockDeleteAvatar.mockResolvedValue(undefined);

      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByRole("button", { name: /削除/ })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole("button", { name: /削除/ }));

      await waitFor(() => {
        expect(mockDeleteAvatar).toHaveBeenCalled();
      });
    });
  });

  // ===== F-10-2: サブスクリプションセクション =====

  describe("サブスクリプションセクション（F-10-2）", () => {
    it("プラン情報が表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
        expect(screen.getByText(/1,980/)).toBeInTheDocument();
      });
    });

    it("ステータスバッジが表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("有効")).toBeInTheDocument();
      });
    });

    it("PORTAL_URL未設定時は管理リンクが非表示", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      expect(screen.queryByRole("link", { name: /管理/ })).not.toBeInTheDocument();
    });
  });

  // ===== F-10-3: メールセクション =====

  describe("メールセクション（F-10-3）", () => {
    it("購読状態が表示される", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByDisplayValue("newsletter@example.com")).toBeInTheDocument();
      });
    });

    it("メールアドレスを新規登録できる", async () => {
      // 購読未登録の状態
      mockGetNewsletterSubscription.mockRejectedValue({
        response: { status: 404 },
      });
      mockSubscribeNewsletter.mockResolvedValue({
        id: 1,
        email: "new@example.com",
        is_active: true,
        created_at: "2026-01-01T00:00:00Z",
        updated_at: "2026-01-01T00:00:00Z",
      });

      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByText("メールアドレス設定")).toBeInTheDocument();
      });

      // メールアドレス入力
      const emailInput = screen.getByPlaceholderText(/メールアドレス/);
      fireEvent.change(emailInput, { target: { value: "new@example.com" } });

      // メールセクション内の保存ボタンをクリック（複数の「保存」ボタンがあるため取得方法を工夫）
      const emailSection = screen.getByText("メールアドレス設定").closest("section")!;
      const saveButton = emailSection.querySelector("button")!;
      fireEvent.click(saveButton);

      await waitFor(() => {
        expect(mockSubscribeNewsletter).toHaveBeenCalledWith({
          email: "new@example.com",
        });
        expect(mockToastSuccess).toHaveBeenCalledWith("購読を開始しました");
      });
    });

    it("メールアドレスを空にして保存すると購読解除される", async () => {
      mockUnsubscribeNewsletter.mockResolvedValue({
        id: 1,
        email: "newsletter@example.com",
        is_active: false,
        created_at: "2026-01-01T00:00:00Z",
        updated_at: "2026-01-01T00:00:00Z",
      });

      render(<SettingsPage />);

      await waitFor(() => {
        expect(screen.getByDisplayValue("newsletter@example.com")).toBeInTheDocument();
      });

      // メールアドレスを空にする
      const emailInput = screen.getByDisplayValue("newsletter@example.com");
      fireEvent.change(emailInput, { target: { value: "" } });

      // メールセクション内の保存ボタンをクリック
      const emailSection = screen.getByText("メールアドレス設定").closest("section")!;
      const saveButton = emailSection.querySelector("button")!;
      fireEvent.click(saveButton);

      await waitFor(() => {
        expect(mockUnsubscribeNewsletter).toHaveBeenCalled();
        expect(mockToastSuccess).toHaveBeenCalledWith("購読を停止しました");
      });
    });
  });

  // ===== F-10-4: FAQセクション =====

  describe("FAQセクション（F-10-4）", () => {
    it("FAQリンクが正しいhrefを持つ", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        const faqLink = screen.getByRole("link", { name: /FAQを見る/ });
        expect(faqLink).toBeInTheDocument();
        expect(faqLink).toHaveAttribute("target", "_blank");
        expect(faqLink).toHaveAttribute("rel", "noopener noreferrer");
      });
    });

    it("お問い合わせリンクが正しいhrefを持つ", async () => {
      render(<SettingsPage />);

      await waitFor(() => {
        const contactLink = screen.getByRole("link", {
          name: /お問い合わせ/,
        });
        expect(contactLink).toBeInTheDocument();
        expect(contactLink).toHaveAttribute("target", "_blank");
        expect(contactLink).toHaveAttribute("rel", "noopener noreferrer");
      });
    });
  });
});
