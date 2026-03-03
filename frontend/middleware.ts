/**
 * Next.js Middleware
 *
 * @description
 * 認証済みユーザーのみがアクセスできるページを保護する。
 *
 * - 未ログイン時: 保護されたページ → /login へリダイレクト
 * - ログイン済み時: /login → / へリダイレクト
 *
 * @see {@link file://./docs/functions/auth/login.md} 認証機能仕様
 */

import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

/**
 * 保護されたページのパス
 */
const protectedPaths = [
  "/",
  "/articles",
  "/article",
  "/news",
  "/poll",
  "/stock",
  "/portfolio",
  "/settings",
];

/**
 * 公開ページのパス（認証不要）
 */
const publicPaths = [
  "/login",
  "/subscription", // 初回ログイン後のサブスクリプション登録画面は一時的に公開
];

/**
 * 認証Cookieの名前
 */
const SESSION_COOKIE_NAME = "session";

function matchesPath(pathname: string, path: string): boolean {
  if (path === "/") {
    return pathname === "/";
  }
  return pathname === path || pathname.startsWith(`${path}/`);
}

function isProtectedPath(pathname: string): boolean {
  return protectedPaths.some((path) => matchesPath(pathname, path));
}

function isPublicPath(pathname: string): boolean {
  return publicPaths.some((path) => matchesPath(pathname, path));
}

/**
 * ミドルウェア関数
 *
 * @param request - Next.jsリクエスト
 * @returns Next.jsレスポンス（リダイレクトまたは継続）
 */
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // セッションCookieの有無を確認
  const sessionCookie = request.cookies.get(SESSION_COOKIE_NAME);
  const hasSession = sessionCookie !== undefined;

  // 保護されたページへのアクセスで未ログインの場合
  if (isProtectedPath(pathname) && !hasSession) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // ログイン済みユーザーが/loginにアクセスした場合（/subscriptionはリダイレクト不要）
  if (pathname === "/login" && hasSession) {
    const url = request.nextUrl.clone();
    url.pathname = "/";
    return NextResponse.redirect(url);
  }

  // その他のページはそのまま
  return NextResponse.next();
}

/**
 * ミドルウェアの適用パス
 *
 * - /api: APIルートは除外（バックエンドで認証）
 * - /_next: Next.js内部ルートは除外
 * - /static: 静的ファイルは除外
 */
export const config = {
  matcher: [
    /*
     * 以下のパスを除外:
     * - /api/* (Next.js API routes)
     * - /_next/* (Next.js internal)
     * - /_next/image (Next.js image optimization)
     * - /favicon.ico, /robots.txt (static files)
     */
    "/((?!api|_next/static|_next/image|favicon.ico|robots.txt).*)",
  ],
};

// テスト用に内部関数を公開
export const __internal__ = {
  matchesPath,
  isProtectedPath,
  isPublicPath,
};
