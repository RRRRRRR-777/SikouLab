/**
 * サブスクリプション登録画面
 *
 * @description
 * 初回ログイン後に表示されるサブスクリプション登録画面。
 * プラン情報を表示し、UnivaPay JSウィジェットでカード決済を行う。
 *
 * @see {@link file://../../../docs/functions/subscription/checkout.md} 詳細設計書
 */

"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import {
  subscriptionApi,
  type Plan,
} from "@/lib/subscription/subscription-api";
import { openCheckoutWidget } from "@/lib/subscription/univapay";
import { authApi } from "@/lib/auth/auth-api";
import { Skeleton } from "@/components/ui/skeleton";

/** 画面の状態 */
type PageState =
  | "loading"
  | "ready"
  | "processing"
  | "polling"
  | "timeout"
  | "error";

/** ポーリング間隔（ミリ秒） */
const POLLING_INTERVAL_MS = 2000;
/** ポーリングタイムアウト（ミリ秒） */
const POLLING_TIMEOUT_MS = 30000;

/**
 * サブスクリプション登録ページコンポーネント
 *
 * 初回ログインユーザーがプラン情報を確認し、
 * UnivaPay JSウィジェットで決済を行う画面。
 *
 * @returns サブスクリプション登録画面
 */
export function SubscriptionPage() {
  const router = useRouter();
  const [plans, setPlans] = useState<Plan[]>([]);
  const [state, setState] = useState<PageState>("loading");
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // プラン情報を取得する
  useEffect(() => {
    subscriptionApi
      .getPlans()
      .then((plans) => {
        setPlans(plans);
        setState("ready");
      })
      .catch(() => {
        setState("error");
      });
  }, []);

  /**
   * subscription_status をポーリングで監視し、active になったらダッシュボードへ遷移する
   */
  const startPolling = useCallback(() => {
    setState("polling");

    const pollInterval = setInterval(async () => {
      try {
        const res = await authApi.getMe();
        if (res.user.subscription_status === "active") {
          clearInterval(pollInterval);
          if (timeoutRef.current) clearTimeout(timeoutRef.current);
          router.push("/");
        }
      } catch {
        // ポーリング中のエラーは無視して継続する
      }
    }, POLLING_INTERVAL_MS);

    pollingRef.current = pollInterval;

    // 30秒タイムアウト
    timeoutRef.current = setTimeout(() => {
      clearInterval(pollInterval);
      toast.error("決済処理がタイムアウトしました");
      setState("timeout");
    }, POLLING_TIMEOUT_MS);
  }, [router]);

  /**
   * UnivaPay ウィジェットを起動し、決済フローを開始する
   */
  const handleCheckout = useCallback(() => {
    const plan = plans[0];
    if (!plan) return;

    setState("processing");

    openCheckoutWidget({
      amount: plan.amount,
      currency: plan.currency,
      onSuccess: async (subscriptionId: string) => {
        try {
          await subscriptionApi.checkout(subscriptionId);
          startPolling();
        } catch (error: unknown) {
          const axiosError = error as {
            response?: { status?: number; data?: { message?: string } };
          };
          if (axiosError.response?.status === 409) {
            toast.error("既に登録済みです");
          } else {
            toast.error("決済処理に失敗しました");
          }
          setState("ready");
        }
      },
      onError: () => {
        toast.error("決済処理に失敗しました");
        setState("ready");
      },
    });
  }, [plans, startPolling]);

  // タイマーのクリーンアップ
  useEffect(() => {
    return () => {
      if (pollingRef.current) clearInterval(pollingRef.current);
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, []);

  const plan = plans[0];
  const isButtonDisabled =
    (state !== "ready" && state !== "timeout") || !plan;

  // ローディング中はスケルトンを表示
  if (state === "loading") {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="w-full max-w-md p-8" data-testid="plan-skeleton">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="mt-2 h-6 w-32" />
          <Skeleton className="mt-4 h-40 w-full" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-md p-8">
        <h1 className="text-2xl font-bold text-foreground">シコウラボ</h1>
        <p className="mt-2 text-foreground">
          ようこそ！サービスをご利用になるにはサブスクリプションの登録が必要です。
        </p>

        {state === "error" && (
          <p className="mt-4 text-destructive">
            プランの取得中にエラーが発生しました
          </p>
        )}

        {plan && (
          <div className="mt-6 rounded-lg border p-6">
            <h2 className="text-xl font-bold text-foreground">{plan.name}</h2>
            <p className="mt-1 text-2xl font-bold text-[#E86D00]">
              ¥{plan.amount.toLocaleString()} / 月
            </p>
            <ul className="mt-4 space-y-2 text-foreground">
              <li>✓ 記事・ニュース閲覧</li>
              <li>✓ 銘柄詳細ページ</li>
              <li>✓ アンケート機能</li>
              <li>✓ ポートフォリオ管理</li>
            </ul>
          </div>
        )}

        {state === "polling" && (
          <p className="mt-4 text-center text-foreground">決済処理中...</p>
        )}

        <button
          type="button"
          onClick={handleCheckout}
          disabled={isButtonDisabled}
          className="mt-6 w-full rounded-lg bg-[#E86D00] px-6 py-4 text-base font-medium text-white transition-all duration-200 hover:bg-[#E86D00]/80 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {state === "timeout"
            ? "もう一度試す"
            : "サブスクリプションを開始する"}
        </button>

        <div className="mt-4 space-y-1 text-center text-sm text-muted-foreground">
          <p>※ いつでもキャンセル可能</p>
        </div>
      </div>
    </div>
  );
}
