/**
 * サブスクリプションAPIクライアント
 *
 * @description
 * プラン取得・サブスクリプション作成APIと通信するクライアント。
 *
 * @module lib/subscription/subscription-api
 */

import { apiClient } from "@/lib/api";

/**
 * プラン情報
 */
export interface Plan {
  /** プランID */
  id: number;
  /** プラン名 */
  name: string;
  /** プラン説明 */
  description: string | null;
  /** 月額料金（最小通貨単位） */
  amount: number;
  /** 通貨コード（ISO-4217） */
  currency: string;
}

/**
 * プラン一覧レスポンス（バックエンドはプラン配列を直接返す）
 */
type PlansResponse = Plan[];

/**
 * チェックアウトレスポンス
 */
interface CheckoutResponse {
  status: string;
}

/**
 * サブスクリプションAPI操作
 */
export const subscriptionApi = {
  /**
   * アクティブなプラン一覧を取得する
   *
   * @returns プラン一覧
   * @throws {AxiosError} API通信エラー時
   */
  async getPlans(): Promise<Plan[]> {
    const response = await apiClient.get<Plan[]>("/plans");
    return response.data ?? [];
  },

  /**
   * UnivaPayサブスクリプションIDをバックエンドに保存する
   *
   * ウィジェットが作成したサブスクリプションのIDをDBに保存する。
   *
   * @param subscriptionId - UnivaPay JSウィジェットが作成したサブスクリプションID
   * @returns チェックアウトステータス（pending）
   * @throws {AxiosError} API通信エラー時（409: 既に登録済み）
   */
  async checkout(subscriptionId: string): Promise<CheckoutResponse> {
    const response = await apiClient.post<CheckoutResponse>(
      "/univapay/checkout",
      {
        subscription_id: subscriptionId,
      },
    );
    return response.data;
  },
};
