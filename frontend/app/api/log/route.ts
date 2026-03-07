/**
 * ログAPIエンドポイント
 *
 * @description
 * クライアントからのログを受け取り、サーバーコンソール（ターミナル）に出力する。
 * LOG_LEVEL 環境変数で指定したレベル以上のログのみ出力する。
 *
 * @route POST /api/log
 */

import { NextRequest, NextResponse } from "next/server";
import type { LogEntry } from "@/lib/logger";

/**
 * ログレベル（数値が小さいほど重要度が高い）
 */
const LOG_LEVEL_PRIORITY: Record<string, number> = {
  error: 0,
  warn: 1,
  info: 2,
  debug: 3,
};

/** 許可するログレベル */
const ALLOWED_LEVELS = ["error", "warn", "info", "debug"] as const;

/** メッセージの最大文字数 */
const MAX_MESSAGE_LENGTH = 1000;

/**
 * クライアントからのログを受け取り、LOG_LEVEL以上のレベルのみサーバーコンソールに出力する。
 *
 * @param request - Next.jsリクエスト
 * @returns 処理結果のJSONレスポンス
 */
export async function POST(request: NextRequest): Promise<NextResponse> {
  try {
    const entry: LogEntry = await request.json();

    // バリデーション: 不正なログレベルやメッセージを拒否
    if (!entry.level || !ALLOWED_LEVELS.includes(entry.level as (typeof ALLOWED_LEVELS)[number])) {
      return NextResponse.json({ error: "Invalid log level" }, { status: 400 });
    }
    if (!entry.message || typeof entry.message !== "string" || entry.message.length > MAX_MESSAGE_LENGTH) {
      return NextResponse.json({ error: "Invalid message" }, { status: 400 });
    }

    // LOG_LEVEL 環境変数で指定したレベル以上のみ出力（デフォルト: info）
    const configuredLevel = process.env.LOG_LEVEL ?? "info";
    const configuredPriority = LOG_LEVEL_PRIORITY[configuredLevel] ?? LOG_LEVEL_PRIORITY.info;
    const entryPriority = LOG_LEVEL_PRIORITY[entry.level] ?? LOG_LEVEL_PRIORITY.info;

    if (entryPriority <= configuredPriority) {
      console.log(JSON.stringify(entry));
    }

    return NextResponse.json({ success: true });
  } catch {
    // パースエラーなどはサイレントで失敗（無限ループ回避）
    return NextResponse.json({ success: false }, { status: 400 });
  }
}
