/**
 * サブスクリプション管理セクション（F-10-2）
 *
 * @description
 * プラン情報表示、ステータスバッジ、管理ボタン（UnivaPayポータルへの遷移）を提供する。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 詳細設計書
 */

"use client";

import { useState, useEffect } from "react";
import { settingsApi } from "@/lib/settings/settings-api";
import type { SubscriptionResponse } from "@/lib/settings/settings-api";
import { Skeleton } from "@/components/ui/skeleton";

const PORTAL_URL = process.env.NEXT_PUBLIC_SUBSCRIPTION_PORTAL_URL ?? "";

/**
 * ステータスの日本語ラベルマッピング
 */
const STATUS_LABELS: Record<string, string> = {
  active: "有効",
  canceled: "解約済み",
  past_due: "支払い遅延",
  trialing: "トライアル",
};

/**
 * ステータスバッジの色マッピング
 */
const STATUS_COLORS: Record<string, string> = {
  active: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200",
  canceled: "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200",
  past_due: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200",
  trialing: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200",
};

/**
 * サブスクリプション管理セクションコンポーネント
 *
 * サブスクリプション状態を表示し、UnivaPayポータルへの遷移ボタンを提供する。
 *
 * @returns サブスクリプション管理セクション
 */
export function SubscriptionSection() {
  const [subscription, setSubscription] = useState<SubscriptionResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // サブスクリプション状態を取得する
  useEffect(() => {
    settingsApi
      .getSubscription()
      .then((data) => {
        setSubscription(data);
      })
      .catch(() => {
        // サブスクリプション未登録の場合はnullのまま
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  return (
    <section className="rounded-lg border border-gray-200 p-4 dark:border-gray-800 lg:p-6">
      <h2 className="text-xl font-bold text-[var(--color-text)]">サブスクリプション</h2>

      {isLoading ? (
        <div className="mt-4 space-y-2">
          <Skeleton className="h-6 w-40" />
          <Skeleton className="h-6 w-24" />
        </div>
      ) : subscription ? (
        <div className="mt-4">
          <p className="text-lg font-bold text-[var(--color-text)]">{subscription.plan_name}</p>
          <div className="mt-1 flex items-center gap-3">
            <span className="text-lg font-medium text-[var(--color-text-secondary)]">
              ¥{subscription.amount.toLocaleString()}/月
            </span>
            <span
              className={`inline-block rounded-full px-2 py-0.5 text-sm font-medium ${STATUS_COLORS[subscription.status] ?? ""}`}
            >
              {STATUS_LABELS[subscription.status] ?? subscription.status}
            </span>
          </div>
          {PORTAL_URL && (
            <div className="mt-4">
              <a
                href={PORTAL_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex min-h-[44px] items-center justify-center rounded-md border border-gray-300 px-4 py-2 text-lg text-[var(--color-text)] hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-800"
              >
                管理
              </a>
            </div>
          )}
        </div>
      ) : (
        <p className="mt-4 text-lg font-medium text-[var(--color-muted-foreground)]">
          サブスクリプションが登録されていません
        </p>
      )}
    </section>
  );
}
