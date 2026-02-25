/**
 * フロントエンドロガー
 *
 * @description
 * クライアント側のログをサーバーコンソール（ターミナル）に出力するためのロガー。
 * バックエンドのzerologと同様に、構造化ログを出力する。
 *
 * @module lib/logger
 */

/**
 * ログレベル
 */
export type LogLevel = "error" | "warn" | "info" | "debug";

/**
 * ログエントリ
 */
export interface LogEntry {
  /** ログレベル */
  level: LogLevel;
  /** メッセージ */
  message: string;
  /** コンテキスト情報 */
  context?: Record<string, unknown>;
  /** タイムスタンプ */
  timestamp: string;
}

/**
 * APIエンドポイント
 */
const LOG_ENDPOINT = "/api/log";

/**
 * サーバーにログを送信する
 *
 * @param entry - ログエントリ
 */
async function sendToServer(entry: LogEntry): Promise<void> {
  try {
    await fetch(LOG_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(entry),
    });
  } catch {
    // ログ送信失敗はサイレントで無視（無限ループ回避）
  }
}

/**
 * ブラウザコンソールに出力する
 *
 * @param entry - ログエントリ
 */
function outputToConsole(entry: LogEntry): void {
  const { level, message, context } = entry;
  const contextStr = context ? ` ${JSON.stringify(context)}` : "";

  switch (level) {
    case "error":
      console.error(`[ERROR] ${message}${contextStr}`);
      break;
    case "warn":
      console.warn(`[WARN] ${message}${contextStr}`);
      break;
    case "info":
      console.info(`[INFO] ${message}${contextStr}`);
      break;
    case "debug":
      console.debug(`[DEBUG] ${message}${contextStr}`);
      break;
  }
}

/**
 * ログを出力する
 *
 * @param level - ログレベル
 * @param message - メッセージ
 * @param context - コンテキスト情報
 */
export function log(level: LogLevel, message: string, context?: Record<string, unknown>): void {
  const entry: LogEntry = {
    level,
    message,
    context,
    timestamp: new Date().toISOString(),
  };

  // サーバーに送信（ターミナルに出力）
  sendToServer(entry);

  // 本番環境ではブラウザコンソールへの出力を抑制する
  if (process.env.NODE_ENV !== "production") {
    outputToConsole(entry);
  }
}

/**
 * エラーログを出力する
 *
 * @param message - メッセージ
 * @param context - コンテキスト情報
 */
export function logError(message: string, context?: Record<string, unknown>): void {
  log("error", message, context);
}

/**
 * 警告ログを出力する
 *
 * @param message - メッセージ
 * @param context - コンテキスト情報
 */
export function logWarn(message: string, context?: Record<string, unknown>): void {
  log("warn", message, context);
}

/**
 * 情報ログを出力する
 *
 * @param message - メッセージ
 * @param context - コンテキスト情報
 */
export function logInfo(message: string, context?: Record<string, unknown>): void {
  log("info", message, context);
}

/**
 * デバッグログを出力する
 *
 * @param message - メッセージ
 * @param context - コンテキスト情報
 */
export function logDebug(message: string, context?: Record<string, unknown>): void {
  log("debug", message, context);
}
