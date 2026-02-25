/**
 * ログイン画面（S-01）
 *
 * @description
 * Firebase OAuth認証によるログイン画面。
 * Google / Apple / X の3つのプロバイダをサポートする。
 *
 * @see {@link file://../../../docs/functions/auth/login.md} 認証機能仕様
 * @see {@link file://../../../docs/versions/1_0_0/SicouLab.pen} Pencilデザイン（h3Lxa）
 */

"use client";

import { useAuth } from "@/lib/auth/auth-context";
import { Loader2 } from "lucide-react";

/**
 * ログイン画面プロパティ
 */
interface LoginPageProps {
  /** OAuthプロバイダからリダイレクトされた場合のエラーメッセージ */
  error?: string;
}

/**
 * OAuthプロバイダの設定
 */
const oauthProviders = [
  {
    id: "google" as const,
    label: "Google でログイン",
    icon: "G",
    bgColor: "bg-[#F5F5F5]",
    hoverColor: "hover:bg-[#E8E8E8]",
    textColor: "text-black",
  },
  {
    id: "apple" as const,
    label: "Apple でログイン",
    icon: "",
    bgColor: "bg-[#F5F5F5]",
    hoverColor: "hover:bg-[#E8E8E8]",
    textColor: "text-black",
  },
  {
    id: "x" as const,
    label: "X でログイン",
    icon: "X",
    bgColor: "bg-[#F5F5F5]",
    hoverColor: "hover:bg-[#E8E8E8]",
    textColor: "text-black",
  },
] as const;

/**
 * ログイン画面コンポーネント
 *
 * OAuthプロバイダ選択ボタンを表示し、ログイン処理を開始する。
 *
 * @param root0 - コンポーネントのプロパティ
 * @param root0.error - OAuthプロバイダからリダイレクトされた場合のエラーメッセージ
 * @returns ログイン画面
 * @example
 * ```tsx
 * <LoginPage />
 * ```
 */
export function LoginPage({ error }: LoginPageProps) {
  const { loginWithGoogle, loginWithApple, loginWithX, isLoading } = useAuth();

  /**
   * OAuthプロバイダ別のログイン処理
   *
   * @param provider - OAuthプロバイダの識別子
   */
  const handleLogin = async (provider: (typeof oauthProviders)[number]["id"]): Promise<void> => {
    try {
      switch (provider) {
        case "google":
          await loginWithGoogle();
          break;
        case "apple":
          await loginWithApple();
          break;
        case "x":
          await loginWithX();
          break;
      }
    } catch {
      // エラーは auth-context でトースト表示済み
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* ブランディングエリア（デスクトップのみ） */}
      <div className="hidden md:flex md:w-1/2 items-center justify-center bg-[#E86D00] p-20">
        <div className="flex flex-col items-center justify-center gap-6 text-center">
          <h1 className="text-5xl font-bold text-white leading-tight">
            ようこそ
            <br />
            シコウラボへ！
          </h1>
          <p className="text-lg text-white opacity-90 max-w-md">
            このコミュニティは米国株投資家が集まり、
            <br />
            専門家の思考プロセスを覗ける場です。
          </p>
        </div>
      </div>

      {/* ログインエリア */}
      <div className="flex w-full flex-col items-center justify-center bg-white p-8 md:w-1/2 dark:bg-black">
        <div className="flex w-full max-w-[400px] flex-col items-center gap-12">
          {/* ロゴ */}
          <div className="flex flex-col items-center gap-4">
            {/* オレンジ色の円（ロゴ） */}
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-[#E86D00]">
              <span className="text-3xl font-bold text-white">S</span>
            </div>
            <h2 className="text-2xl font-bold text-black dark:text-white">SikouLab</h2>
          </div>

          {/* OAuthボタン */}
          <div className="flex w-full flex-col gap-4">
            {oauthProviders.map((provider) => (
              <button
                key={provider.id}
                type="button"
                onClick={() => handleLogin(provider.id)}
                disabled={isLoading}
                className={`
                  flex h-14 w-full items-center justify-center gap-4 rounded-lg
                  ${provider.bgColor} ${provider.hoverColor}
                  ${provider.textColor}
                  px-6 py-4 text-base font-medium
                  transition-all duration-200
                  disabled:cursor-not-allowed disabled:opacity-50
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#E86D00] focus-visible:ring-offset-2
                `}
              >
                {isLoading ? (
                  <Loader2 className="h-5 w-5 animate-spin" />
                ) : provider.icon ? (
                  <span className="text-lg font-bold">{provider.icon}</span>
                ) : (
                  <svg
                    className="h-5 w-5"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    aria-hidden="true"
                  >
                    <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z" />
                  </svg>
                )}
                <span>{provider.label}</span>
              </button>
            ))}
          </div>

          {/* 利用規約への同意文言 */}
          <p className="text-center text-xs text-gray-500">
            ログインすることで利用規約に同意したものとみなされます
          </p>

          {/* エラーメッセージ（OAuthプロバイダからリダイレクトされた場合） */}
          {error && (
            <div className="w-full rounded-md bg-destructive/15 p-3 text-sm text-destructive">
              {error}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
