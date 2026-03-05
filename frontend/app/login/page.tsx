/**
 * ログインページ
 *
 * @description
 * 未ログインユーザーがアクセスするログイン画面。
 * OAuth認証によるログインを提供する。
 *
 * @route /login
 */

import { LoginPage } from "@/components/auth/LoginPage";

/**
 * ログインページコンポーネント
 *
 * URLパラメータからエラーメッセージを受け取る場合がある。
 *
 * @param root0 - ページのプロパティ
 * @param root0.searchParams - URLクエリパラメータ
 * @returns ログインページコンポーネント
 * @example
 * - /login - 通常のログイン画面
 * - /login?error=auth_failed - エラーメッセージ付き
 */
export default async function Login({
  searchParams,
}: {
  searchParams: Promise<{ error?: string }>;
}) {
  const { error } = await searchParams;
  return <LoginPage error={error} />;
}
