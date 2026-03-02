/**
 * サブスクリプション登録ページ
 *
 * @description
 * 初回ログイン後にリダイレクトされるサブスクリプション登録画面。
 * UnivaPay JSウィジェットで決済を行い、ダッシュボードへ遷移する。
 *
 * @route /subscription
 */

import { SubscriptionPage } from "@/components/subscription/SubscriptionPage";

/**
 * サブスクリプション登録ページコンポーネント
 *
 * @returns サブスクリプション登録画面
 */
export default function Subscription() {
  return <SubscriptionPage />;
}
