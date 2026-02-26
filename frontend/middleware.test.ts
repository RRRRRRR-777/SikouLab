import { describe, expect, it, vi } from "vitest";　
import { NextRequest } from "next/server";

import { __internal__, middleware } from "./middleware";

describe("middleware path matching", () => {
  it("treats '/' as exact match only", () => {
    expect(__internal__.matchesPath("/", "/")).toBe(true);
    expect(__internal__.matchesPath("/login", "/")).toBe(false);
    expect(__internal__.matchesPath("/articles", "/")).toBe(false);
  });

  it("matches exact path and sub paths for non-root routes", () => {
    expect(__internal__.matchesPath("/articles", "/articles")).toBe(true);
    expect(__internal__.matchesPath("/articles/123", "/articles")).toBe(true);
    expect(__internal__.matchesPath("/articles-legacy", "/articles")).toBe(false);
  });

  it("keeps /login as public and protected routes as expected", () => {
    expect(__internal__.isPublicPath("/login")).toBe(true);
    expect(__internal__.isPublicPath("/subscription")).toBe(true);
    expect(__internal__.isProtectedPath("/login")).toBe(false);
    expect(__internal__.isProtectedPath("/")).toBe(true);
    expect(__internal__.isProtectedPath("/articles/2026")).toBe(true);
  });
});

/**
 * ミドルウェアのルートガードテスト
 *
 * @description
 * middleware 関数本体の動作を確認する。
 * セッションCookieの有無に応じたリダイレクト・通過の挙動をテストする。
 */
describe("middleware route guard", () => {
  /**
   * テスト用のNextRequestを生成するヘルパー
   *
   * @param pathname - リクエスト先のパス
   * @param hasSession - セッションCookieを付与するかどうか
   * @returns NextRequest インスタンス
   */
  function makeRequest(pathname: string, hasSession: boolean): NextRequest {
    const url = `http://localhost${pathname}`;
    const req = new NextRequest(url);
    if (hasSession) {
      req.cookies.set("session", "dummy-session-token");
    }
    return req;
  }

  it("セッションなしで保護パスへのリクエストは/loginにリダイレクト", () => {
    // 準備: セッションなしで保護パス（/）へのリクエスト
    const req = makeRequest("/", false);

    // 実行
    const res = middleware(req);

    // 検証: /login へリダイレクト
    expect(res.status).toBe(307);
    expect(res.headers.get("location")).toContain("/login");
  });

  it("セッションありで/loginへのリクエストは/にリダイレクト", () => {
    // 準備: セッションありで /login へのリクエスト
    const req = makeRequest("/login", true);

    // 実行
    const res = middleware(req);

    // 検証: / へリダイレクト
    expect(res.status).toBe(307);
    expect(res.headers.get("location")).toContain("/");
    expect(res.headers.get("location")).not.toContain("/login");
  });

  it("セッションなしで/subscriptionへのアクセスは通過する", () => {
    // 準備: セッションなしで公開パス（/subscription）へのリクエスト
    const req = makeRequest("/subscription", false);

    // 実行
    const res = middleware(req);

    // 検証: リダイレクトせずにそのまま通過（NextResponse.next()）
    expect(res.status).toBe(200);
    expect(res.headers.get("location")).toBeNull();
  });

  it("セッションありで保護パスへのアクセスは通過する", () => {
    // 準備: セッションありで保護パス（/articles）へのリクエスト
    const req = makeRequest("/articles", true);

    // 実行
    const res = middleware(req);

    // 検証: リダイレクトせずにそのまま通過（NextResponse.next()）
    expect(res.status).toBe(200);
    expect(res.headers.get("location")).toBeNull();
  });
});
