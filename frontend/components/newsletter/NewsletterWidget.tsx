/**
 * メール登録ウィジェット
 *
 * @description
 * 各画面のページ下部に配置する再利用可能なニュースレター購読コンポーネント。
 * ログイン済みユーザーにのみ表示され、既に購読済みの場合は購読済み表示になる。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 設定画面仕様
 * @see {@link file://../../../docs/versions/1_0_0/SicouLab.pen} Pencilデザイン
 */

"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useAuth } from "@/lib/auth/auth-context";
import { settingsApi } from "@/lib/settings/settings-api";
import { EMAIL_REGEX } from "@/lib/validation";

/** ウィジェット外枠の共通スタイル */
const CONTAINER_CLASS =
  "w-full rounded-lg border border-gray-200 bg-[var(--color-bg)] p-6 dark:border-gray-800";

/**
 * メール登録ウィジェットコンポーネント
 *
 * ログイン済みユーザーにニュースレター購読登録フォームを表示する。
 * 未ログイン時や認証ローディング中は何も表示しない。
 *
 * @returns ニュースレター購読ウィジェット（未ログイン時はnull）
 *
 * @example
 * ```tsx
 * <NewsletterWidget />
 * ```
 */
export function NewsletterWidget() {
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubscribed, setIsSubscribed] = useState(false);
  const [isCheckingSubscription, setIsCheckingSubscription] = useState(true);

  // ログイン済みの場合、購読状況を確認する
  useEffect(() => {
    if (!isAuthenticated) {
      setIsCheckingSubscription(false);
      return;
    }

    const checkSubscription = async () => {
      try {
        const subscription = await settingsApi.getNewsletterSubscription();
        if (subscription.is_active) {
          setIsSubscribed(true);
        }
      } catch {
        // 404（購読未登録）は正常系として扱う
      } finally {
        setIsCheckingSubscription(false);
      }
    };

    checkSubscription();
  }, [isAuthenticated]);

  // 未ログイン・認証ローディング中は表示しない
  if (!isAuthenticated || authLoading || !user) {
    return null;
  }

  // 購読状況チェック中は表示しない
  if (isCheckingSubscription) {
    return null;
  }

  /**
   * メールアドレスのクライアント側バリデーションを実行する。
   *
   * @param value - バリデーション対象のメールアドレス
   * @returns バリデーション結果（true: 有効、false: 無効）
   */
  const validateEmail = (value: string): boolean => {
    if (!value.trim()) {
      setError("メールアドレスを入力してください");
      return false;
    }
    if (!EMAIL_REGEX.test(value)) {
      setError("有効なメールアドレスを入力してください");
      return false;
    }
    setError("");
    return true;
  };

  /**
   * 購読登録フォームの送信ハンドラ。
   *
   * バリデーション → API呼び出し → トースト表示の順に処理する。
   */
  const handleSubmit = async () => {
    if (!validateEmail(email)) {
      return;
    }

    setIsSubmitting(true);
    try {
      await settingsApi.subscribeNewsletter({ email });
      toast.success("購読を開始しました");
      setIsSubscribed(true);
      setEmail("");
      setError("");
    } catch (err: unknown) {
      // バックエンドのエラーメッセージがある場合はそれを表示する
      const apiError = err as {
        response?: { data?: { message?: string } };
      };
      if (apiError.response?.data?.message) {
        toast.error(apiError.response.data.message);
      } else {
        toast.error("購読登録に失敗しました");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  // 購読済み表示
  if (isSubscribed) {
    return (
      <div className={CONTAINER_CLASS}>
        <h3 className="text-lg font-bold text-[var(--color-text)]">
          メールで記事を受け取る
        </h3>
        <p className="mt-2 text-base text-[#E86D00] font-medium">
          購読済みです
        </p>
      </div>
    );
  }

  // 購読登録フォーム
  return (
    <div className={CONTAINER_CLASS}>
      <h3 className="text-lg font-bold text-[var(--color-text)]">
        メールで記事を受け取る
      </h3>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        毎朝7:30に記事の要約をお届けします。
      </p>
      <div className="mt-4">
        <input
          type="email"
          placeholder="メールアドレス"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            // 入力中にエラーをクリアする
            if (error) setError("");
          }}
          className="w-full rounded-md border border-gray-200 bg-[var(--color-bg)] px-4 py-3 text-base text-[var(--color-text)] placeholder:text-[var(--color-muted)] focus:border-[#E86D00] focus:outline-none focus:ring-1 focus:ring-[#E86D00] dark:border-gray-800"
        />
        {error && (
          <p className="mt-1 text-sm text-red-500">{error}</p>
        )}
      </div>
      <button
        type="button"
        onClick={handleSubmit}
        disabled={isSubmitting}
        className="mt-3 min-h-[44px] w-full rounded-md bg-[#E86D00] px-4 py-3 text-base font-medium text-white transition-colors hover:bg-[#D05E00] disabled:cursor-not-allowed disabled:opacity-50"
      >
        購読する
      </button>
    </div>
  );
}
