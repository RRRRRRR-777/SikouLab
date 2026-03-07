/**
 * UnivaPay JSウィジェットラッパー
 *
 * @description
 * UnivaPay JS SDKのチェックアウトウィジェットを起動するラッパー関数。
 * テスト時はモックに差し替えて使用する。
 *
 * @see https://docs.univapay.com/en/widget-javascript/
 * @module lib/subscription/univapay
 */

/** UnivaPay checkout.js がグローバルに公開する型 */
interface UnivapayCheckoutInstance {
  open: () => void;
}

interface UnivapayCheckoutStatic {
  create: (params: Record<string, unknown>) => UnivapayCheckoutInstance;
}

declare global {
  interface Window {
    UnivapayCheckout?: UnivapayCheckoutStatic;
  }
}

/**
 * チェックアウトウィジェットの設定
 */
export interface CheckoutWidgetConfig {
  /** 月額料金（最小通貨単位） */
  amount: number;
  /** 通貨コード（ISO-4217） */
  currency: string;
  /** 決済成功コールバック（subscription_idを受け取る） */
  onSuccess: (subscriptionId: string) => void;
  /** 決済エラーコールバック */
  onError?: (error: Error) => void;
  /** ウィジェットがキャンセル・クローズされた時のコールバック */
  onClose?: () => void;
}

/** checkout.js の CDN URL */
const CHECKOUT_SCRIPT_URL = "https://widget.univapay.com/client/checkout.js";

/**
 * checkout.js スクリプトを動的にロードする
 *
 * 既にロード済みの場合はスキップする。
 *
 * @returns スクリプトロード完了のPromise
 */
function loadCheckoutScript(): Promise<void> {
  if (window.UnivapayCheckout) {
    return Promise.resolve();
  }

  return new Promise((resolve, reject) => {
    // 既に <script> タグが存在するがまだロード中の場合
    const existing = document.querySelector(
      `script[src="${CHECKOUT_SCRIPT_URL}"]`,
    );
    if (existing) {
      existing.addEventListener("load", () => resolve());
      existing.addEventListener("error", () =>
        reject(new Error("UnivaPay SDK の読み込みに失敗しました")),
      );
      return;
    }

    const script = document.createElement("script");
    script.src = CHECKOUT_SCRIPT_URL;
    script.async = true;
    script.onload = () => resolve();
    script.onerror = () =>
      reject(new Error("UnivaPay SDK の読み込みに失敗しました"));
    document.head.appendChild(script);
  });
}

/**
 * UnivaPay JSウィジェットを起動する
 *
 * checkout.js を動的にロードし、UnivapayCheckout.create() でウィジェットを生成して開く。
 * token-created イベントで transaction_token_id を取得し、onSuccess コールバックを呼ぶ。
 *
 * @param config - ウィジェット設定（金額・通貨・コールバック）
 * @returns ウィジェット起動のPromise
 */
export async function openCheckoutWidget(
  config: CheckoutWidgetConfig,
): Promise<void> {
  const appId = process.env.NEXT_PUBLIC_UNIVAPAY_APP_ID;
  if (!appId) {
    config.onError?.(
      new Error("NEXT_PUBLIC_UNIVAPAY_APP_ID が設定されていません"),
    );
    return;
  }

  try {
    await loadCheckoutScript();
  } catch (error) {
    config.onError?.(error as Error);
    return;
  }

  if (!window.UnivapayCheckout) {
    config.onError?.(new Error("UnivaPay SDK の初期化に失敗しました"));
    return;
  }

  // checkout: "payment" でウィジェットがサブスクリプション作成と3Dセキュア認証を一括処理する。
  // checkout: "token" ではトークン作成のみで、サーバーサイドでサブスクリプション作成時に
  // 3Dセキュアが発生するがユーザーが認証操作できないため課金失敗する。
  const checkout = window.UnivapayCheckout.create({
    appId,
    checkout: "payment",
    amount: config.amount,
    currency: config.currency,
    tokenType: "subscription",
    subscriptionPeriod: "monthly",
  });

  // サブスクリプション作成成功イベント（3Dセキュア認証完了後に発火）
  window.addEventListener(
    "univapay:subscription-created",
    ((event: CustomEvent<{ id: string }>) => {
      config.onSuccess(event.detail.id);
    }) as EventListener,
    { once: true },
  );

  // エラーイベント
  window.addEventListener(
    "univapay:error",
    ((event: CustomEvent<{ message: string }>) => {
      config.onError?.(new Error(event.detail?.message ?? "決済エラー"));
    }) as EventListener,
    { once: true },
  );

  // ウィジェットクローズイベント（キャンセル時）
  window.addEventListener(
    "univapay:closed",
    (() => {
      config.onClose?.();
    }) as EventListener,
    { once: true },
  );

  checkout.open();
}
