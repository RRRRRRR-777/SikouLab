/**
 * グローバルエラーハンドラ
 *
 * window.onerror と unhandledrejection をフックし、
 * クライアントサイドの未捕捉エラーを既存の logger 経由でサーバーログに送信する。
 *
 * @module components/GlobalErrorHandler
 */

"use client";

import { useEffect } from "react";
import { logError } from "@/lib/logger";

/**
 * グローバルエラーハンドラコンポーネント
 *
 * ルートレイアウトに配置することで、ブラウザの未捕捉エラーを
 * サーバーログ（/api/log 経由）に送信する。
 *
 * @returns null（UIを持たない）
 */
export function GlobalErrorHandler() {
  useEffect(() => {
    const handleError = (event: ErrorEvent) => {
      logError("未捕捉エラー", {
        message: event.message,
        source: event.filename,
        line: event.lineno,
        column: event.colno,
        stack: event.error?.stack,
      });
    };

    const handleRejection = (event: PromiseRejectionEvent) => {
      const reason = event.reason;
      logError("未処理のPromise拒否", {
        message: reason?.message ?? String(reason),
        stack: reason?.stack,
      });
    };

    window.addEventListener("error", handleError);
    window.addEventListener("unhandledrejection", handleRejection);

    return () => {
      window.removeEventListener("error", handleError);
      window.removeEventListener("unhandledrejection", handleRejection);
    };
  }, []);

  return null;
}
