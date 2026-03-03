/**
 * サブスクリプション登録画面のテスト
 *
 * @description
 * 機能要件4（サブスクリプション登録画面）の単体テスト。
 * プラン表示・UnivaPay ウィジェット起動・ポーリング・タイムアウトをテストする。
 *
 * @see {@link file://../../../../docs/functions/subscription/checkout.md} 詳細設計書
 */

import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor, fireEvent, act } from "@testing-library/react";
import { AxiosError } from "axios";
import { SubscriptionPage } from "../SubscriptionPage";

// --- モック変数のホイスティング ---
const {
  mockPush,
  mockToastError,
  mockToastSuccess,
  mockGetPlans,
  mockCheckout,
  mockGetMe,
  mockOpenWidget,
} = vi.hoisted(() => ({
  mockPush: vi.fn(),
  mockToastError: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockGetPlans: vi.fn(),
  mockCheckout: vi.fn(),
  mockGetMe: vi.fn(),
  mockOpenWidget: vi.fn(),
}));

// --- 依存モジュールのモック ---

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
  usePathname: () => "/subscription",
}));

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

vi.mock("@/lib/subscription/subscription-api", () => ({
  subscriptionApi: {
    getPlans: (...args: unknown[]) => mockGetPlans(...args),
    checkout: (...args: unknown[]) => mockCheckout(...args),
  },
}));

vi.mock("@/lib/auth/auth-api", () => ({
  authApi: {
    getMe: (...args: unknown[]) => mockGetMe(...args),
  },
}));

vi.mock("@/lib/subscription/univapay", () => ({
  openCheckoutWidget: (...args: unknown[]) => mockOpenWidget(...args),
}));

// --- テストデータ ---

const mockPlan = {
  id: 1,
  name: "プレミアムプラン",
  description: "全機能利用可能",
  amount: 1980,
  currency: "JPY",
};

// --- テスト ---

describe("SubscriptionPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // ===== 機能仕様1: プラン情報を表示する =====

  describe("プラン表示（機能要件4/機能仕様1）", () => {
    it("ローディング中はスケルトンを表示する", () => {
      // API レスポンス未到着の状態を再現
      mockGetPlans.mockReturnValue(new Promise(() => {}));

      render(<SubscriptionPage />);

      expect(screen.getByTestId("plan-skeleton")).toBeInTheDocument();
    });

    it("プラン情報が正しく表示される", async () => {
      mockGetPlans.mockResolvedValue([mockPlan]);

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
        expect(screen.getByText(/¥1,980/)).toBeInTheDocument();
      });
    });

    it("プラン取得APIエラー時にエラー表示", async () => {
      mockGetPlans.mockRejectedValue(new Error("Server Error"));

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText(/エラー/)).toBeInTheDocument();
      });
      // ボタンが非活性
      expect(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      ).toBeDisabled();
    });

    it("プランが0件の場合はボタンが非活性", async () => {
      mockGetPlans.mockResolvedValue([]);

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
        ).toBeDisabled();
      });
    });
  });

  // ===== 機能仕様2: UnivaPay JS ウィジェットを起動する =====

  describe("ウィジェット起動（機能要件4/機能仕様2）", () => {
    it("ボタンクリックでウィジェットを初期化する", async () => {
      mockGetPlans.mockResolvedValue([mockPlan]);

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      expect(mockOpenWidget).toHaveBeenCalledWith(
        expect.objectContaining({
          amount: 1980,
          currency: "JPY",
        }),
      );
    });

    it("プラン未取得時はボタン非活性", () => {
      // API レスポンス待ち中
      mockGetPlans.mockReturnValue(new Promise(() => {}));

      render(<SubscriptionPage />);

      const button = screen.queryByRole("button", {
        name: /サブスクリプションを開始する/,
      });
      // ボタンがまだレンダリングされていないか、disabled の状態
      if (button) {
        expect(button).toBeDisabled();
      }
    });

    it("処理中はボタンが非活性になる", async () => {
      mockGetPlans.mockResolvedValue([mockPlan]);
      // ウィジェットのコールバックを即座に呼ぶ
      mockOpenWidget.mockImplementation(
        ({ onSuccess }: { onSuccess: (tokenId: string) => void }) => {
          onSuccess("tok_xxx");
        },
      );
      // checkout が pending 状態で止まる
      mockCheckout.mockReturnValue(new Promise(() => {}));

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
        ).toBeDisabled();
      });
    });
  });

  // ===== 機能仕様3: 決済完了後にダッシュボードへ遷移する =====

  describe("決済完了後フロー（機能要件4/機能仕様3）", () => {
    /**
     * ウィジェット→checkout→ポーリング開始までの共通セットアップ
     */
    const setupCheckoutFlow = () => {
      mockGetPlans.mockResolvedValue([mockPlan]);
      mockOpenWidget.mockImplementation(
        ({ onSuccess }: { onSuccess: (tokenId: string) => void }) => {
          onSuccess("tok_xxx");
        },
      );
    };

    it("トークンコールバック後にcheckout APIを呼ぶ", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe.mockResolvedValue({ user: { subscription_status: null } });

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      await waitFor(() => {
        expect(mockCheckout).toHaveBeenCalledWith("tok_xxx");
      });
    });

    it("checkout成功後に「決済処理中...」を表示", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe.mockResolvedValue({ user: { subscription_status: null } });

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      await waitFor(() => {
        expect(screen.getByText(/決済処理中/)).toBeInTheDocument();
      });
    });

    it("checkout APIエラー（500）でエラートースト", async () => {
      setupCheckoutFlow();
      mockCheckout.mockRejectedValue({ response: { status: 500 } });

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      await waitFor(() => {
        expect(mockToastError).toHaveBeenCalled();
      });

      // ウィジェット前の状態に戻る（ボタンが再度有効）
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
        ).not.toBeDisabled();
      });
    });

    it("既にactive（409）でエラートースト", async () => {
      setupCheckoutFlow();
      const axiosError = new AxiosError("Request failed with status code 409");
      axiosError.response = { status: 409, data: { message: "既に登録済みです" } } as import("axios").AxiosResponse;
      mockCheckout.mockRejectedValue(axiosError);

      render(<SubscriptionPage />);

      await waitFor(() => {
        expect(screen.getByText("プレミアムプラン")).toBeInTheDocument();
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      await waitFor(() => {
        expect(mockToastError).toHaveBeenCalledWith(
          expect.stringContaining("登録済み"),
        );
      });
    });
  });

  // ===== 機能仕様3: ポーリング・タイムアウト（fake timers使用） =====

  describe("ポーリング・タイムアウト（機能要件4/機能仕様3）", () => {
    /**
     * ウィジェット→checkout→ポーリング開始までの共通セットアップ（fake timers版）
     */
    const setupCheckoutFlow = () => {
      mockGetPlans.mockResolvedValue([mockPlan]);
      mockOpenWidget.mockImplementation(
        ({ onSuccess }: { onSuccess: (tokenId: string) => void }) => {
          onSuccess("tok_xxx");
        },
      );
    };

    /**
     * checkout完了までの共通フロー（fake timers環境下）
     * waitFor が使えないため act + advanceTimersByTimeAsync で同期する
     */
    const renderAndStartCheckout = async () => {
      render(<SubscriptionPage />);

      // useEffect（getPlans）のPromise解決を待つ
      await act(async () => {
        await vi.advanceTimersByTimeAsync(0);
      });

      fireEvent.click(
        screen.getByRole("button", { name: /サブスクリプションを開始する/ }),
      );

      // checkout API のPromise解決 → startPolling
      await act(async () => {
        await vi.advanceTimersByTimeAsync(0);
      });
    };

    beforeEach(() => {
      vi.useFakeTimers();
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it("ポーリングでactive検知後にダッシュボードへ遷移", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe
        .mockResolvedValueOnce({ user: { subscription_status: null } })
        .mockResolvedValueOnce({ user: { subscription_status: "active" } });

      await renderAndStartCheckout();

      // ポーリング1回目（null）
      await act(async () => {
        await vi.advanceTimersByTimeAsync(2000);
      });

      // ポーリング2回目（active）
      await act(async () => {
        await vi.advanceTimersByTimeAsync(2000);
      });

      expect(mockPush).toHaveBeenCalledWith("/");
    });

    it("ポーリング間隔が2秒", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe.mockResolvedValue({ user: { subscription_status: null } });

      await renderAndStartCheckout();

      const callCountBefore = mockGetMe.mock.calls.length;

      // 1秒後にはまだ呼ばれていない
      await act(async () => {
        await vi.advanceTimersByTimeAsync(1000);
      });
      expect(mockGetMe.mock.calls.length).toBe(callCountBefore);

      // さらに1秒後（合計2秒）にポーリングが1回実行される
      await act(async () => {
        await vi.advanceTimersByTimeAsync(1000);
      });
      expect(mockGetMe.mock.calls.length).toBeGreaterThan(callCountBefore);
    });

    it("30秒タイムアウトでエラートースト表示", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe.mockResolvedValue({ user: { subscription_status: null } });

      await renderAndStartCheckout();

      // 30秒経過
      await act(async () => {
        await vi.advanceTimersByTimeAsync(30000);
      });

      expect(mockToastError).toHaveBeenCalled();
    });

    it("タイムアウト後に「もう一度試す」ボタンが表示される", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe.mockResolvedValue({ user: { subscription_status: null } });

      await renderAndStartCheckout();

      await act(async () => {
        await vi.advanceTimersByTimeAsync(30000);
      });

      expect(
        screen.getByRole("button", { name: /もう一度試す/ }),
      ).toBeInTheDocument();
    });

    it("ポーリング中のAPIエラーでポーリング継続", async () => {
      setupCheckoutFlow();
      mockCheckout.mockResolvedValue({ status: "pending" });
      mockGetMe
        .mockRejectedValueOnce(new Error("Server Error"))
        .mockResolvedValueOnce({ user: { subscription_status: "active" } });

      await renderAndStartCheckout();

      // 1回目のポーリング（エラー → 中断しない）
      await act(async () => {
        await vi.advanceTimersByTimeAsync(2000);
      });

      // 2回目のポーリング（active）
      await act(async () => {
        await vi.advanceTimersByTimeAsync(2000);
      });

      expect(mockPush).toHaveBeenCalledWith("/");
    });
  });
});
