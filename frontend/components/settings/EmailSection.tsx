/**
 * メールアドレス設定セクション（F-10-3）
 *
 * @description
 * ログイン用メールアドレス表示（読み取り専用）、
 * ニュースレター用メール登録/変更を提供する。
 * メールアドレスを入力して保存するだけのシンプルなUI。
 * メールアドレスが空の状態で保存すると購読解除になる。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 詳細設計書
 */

"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useAuth } from "@/lib/auth/auth-context";
import { settingsApi } from "@/lib/settings/settings-api";
import type { NewsletterSubscriptionResponse } from "@/lib/settings/settings-api";
import { Skeleton } from "@/components/ui/skeleton";
import { EMAIL_REGEX } from "@/lib/validation";

/**
 * メールアドレス設定セクションコンポーネント
 *
 * ログイン用メールの表示とニュースレター用メールアドレスの管理UIを提供する。
 * メールアドレスを入力して保存するシンプルなUI。空にして保存で購読解除。
 *
 * @returns メールアドレス設定セクション
 */
export function EmailSection() {
  const { user } = useAuth();
  const [newsletterEmail, setNewsletterEmail] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  // 購読レコードが存在するかどうか。サーバー側にレコードがあればtrue。
  const [hasExistingSubscription, setHasExistingSubscription] = useState(false);

  // 購読状態を取得する
  useEffect(() => {
    settingsApi
      .getNewsletterSubscription()
      .then((data: NewsletterSubscriptionResponse) => {
        setNewsletterEmail(data.email ?? "");
        setHasExistingSubscription(true);
      })
      .catch(() => {
        // 404: 購読未登録の場合はデフォルト状態のまま
        setHasExistingSubscription(false);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  /**
   * メール設定を保存する。
   *
   * メールアドレスが入力されていれば購読登録/更新、空なら購読解除する。
   */
  const handleSave = async () => {
    const trimmedEmail = newsletterEmail.trim();

    // メールアドレスが入力されている場合はバリデーション
    if (trimmedEmail && !EMAIL_REGEX.test(trimmedEmail)) {
      toast.error("有効なメールアドレスを入力してください");
      return;
    }

    setIsSaving(true);
    try {
      if (trimmedEmail) {
        if (hasExistingSubscription) {
          // 既存購読のメールアドレス変更
          await settingsApi.updateNewsletterEmail({ email: trimmedEmail });
          toast.success("メールアドレスを保存しました");
        } else {
          // 新規購読登録
          await settingsApi.subscribeNewsletter({ email: trimmedEmail });
          setHasExistingSubscription(true);
          toast.success("購読を開始しました");
        }
      } else {
        if (hasExistingSubscription) {
          // メールアドレスが空 → 購読解除
          await settingsApi.unsubscribeNewsletter();
          setHasExistingSubscription(false);
          toast.success("購読を停止しました");
        }
      }
    } catch {
      toast.error("保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <section className="rounded-lg border border-gray-200 p-4 dark:border-gray-800 lg:p-6">
      <h2 className="text-xl font-bold text-[var(--color-text)]">メールアドレス設定</h2>

      {isLoading ? (
        <div className="mt-4 space-y-2">
          <Skeleton className="h-6 w-60" />
          <Skeleton className="h-10 w-full" />
        </div>
      ) : (
        <div className="mt-4 space-y-4">
          {/* ログイン用メールアドレス（読み取り専用） */}
          {user?.oauthProvider && (
            <div>
              <label className="text-lg font-bold text-[var(--color-text)]">
                ログイン用メールアドレス
              </label>
              <p className="mt-1 text-lg font-medium text-[var(--color-text-secondary)]">
                OAuthプロバイダから取得
              </p>
            </div>
          )}

          {/* ニュースレター用メールアドレス */}
          <div>
            <label
              htmlFor="newsletter-email"
              className="text-lg font-bold text-[var(--color-text)]"
            >
              ニュースレター用メールアドレス
            </label>
            <p className="mt-1 text-sm text-[var(--color-muted-foreground)]">
              毎朝7:30に記事の要約をお届けします。空にして保存すると購読を停止します。
            </p>
            <div className="mt-2 flex gap-2">
              <input
                id="newsletter-email"
                type="email"
                value={newsletterEmail}
                onChange={(e) => setNewsletterEmail(e.target.value)}
                placeholder="メールアドレスを入力"
                className="min-h-[44px] flex-1 rounded-md border border-gray-300 bg-[var(--color-bg)] px-3 py-2 text-lg text-[var(--color-text)] outline-none focus:ring-2 focus:ring-[var(--color-primary)] dark:border-gray-600"
              />
              <button
                type="button"
                onClick={handleSave}
                disabled={isSaving}
                className="min-h-[44px] rounded-md bg-[var(--color-primary)] px-4 py-2 text-lg text-white hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-50"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}
