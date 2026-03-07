/**
 * 設定画面のAPIクライアント
 *
 * @description
 * プロフィール更新、サブスクリプション状態取得、ニュースレター管理APIと通信する。
 * openapi.yaml の定義に準拠したリクエスト/レスポンス型を使用する。
 *
 * @module lib/settings/settings-api
 */

import axios from "axios";
import { apiClient } from "@/lib/api";
import type { User } from "@/lib/auth/types";

// --- 型定義 ---

/**
 * プロフィール更新リクエスト
 */
export interface UpdateProfileRequest {
  /** 表示名（1-50文字、空白のみ不可） */
  display_name: string;
}

/**
 * ユーザー情報レスポンス（バックエンドのUserResponseスキーマ）
 */
export interface UserResponse {
  user: User;
}

/**
 * アバターアップロードレスポンス
 */
export interface AvatarResponse {
  /** アップロードされたアバター画像のURL */
  avatar_url: string;
}

/**
 * サブスクリプション状態レスポンス
 */
export interface SubscriptionResponse {
  /** プラン名 */
  plan_name: string;
  /** 月額料金（最小通貨単位） */
  amount: number;
  /** 通貨コード（ISO-4217） */
  currency: string;
  /** サブスクリプション状態 */
  status: "active" | "canceled" | "past_due" | "trialing";
}

/**
 * ポータルURLレスポンス
 */
export interface PortalResponse {
  /** UnivaPayカスタマーポータルURL */
  portal_url: string;
}

/**
 * ニュースレター購読状態レスポンス
 */
export interface NewsletterSubscriptionResponse {
  /** 購読ID */
  id: number;
  /** ニュースレター用メールアドレス */
  email: string;
  /** 購読状態（true=購読中、false=停止中） */
  is_active: boolean;
  /** 作成日時 */
  created_at: string;
  /** 更新日時 */
  updated_at: string;
}

/**
 * ニュースレター購読登録リクエスト
 */
export interface NewsletterSubscribeRequest {
  /** ニュースレター用メールアドレス */
  email: string;
}

/**
 * ニュースレター メールアドレス変更リクエスト
 */
export interface NewsletterUpdateRequest {
  /** 新しいニュースレター用メールアドレス */
  email: string;
}

// --- API操作 ---

/**
 * 設定画面のAPI操作をまとめたオブジェクト。
 *
 * 各メソッドは openapi.yaml に定義されたエンドポイントと1対1で対応する。
 */
export const settingsApi = {
  /**
   * 表示名を更新する。
   *
   * @param data - 更新リクエスト（display_name: 1-50文字、空白のみ不可）
   * @returns 更新後のユーザー情報
   * @throws {AxiosError} 400: バリデーションエラー, 401: 未認証
   */
  async updateProfile(data: UpdateProfileRequest): Promise<UserResponse> {
    const response = await apiClient.patch<UserResponse>("/users/me", data);
    return response.data;
  },

  /**
   * アバター画像をアップロードする。
   *
   * @param file - アバター画像ファイル（JPEG/PNG/GIF、最大5MB）
   * @returns アップロードされた画像のURL
   * @throws {AxiosError} 400: 不正なファイル形式/サイズ超過, 401: 未認証
   */
  async uploadAvatar(file: File): Promise<AvatarResponse> {
    const formData = new FormData();
    formData.append("image", file);
    // apiClientのデフォルトContent-Type（application/json）がFormDataを壊すため、
    // ヘッダーなしの素のaxiosで送信しboundaryを自動設定させる。
    const baseURL = apiClient.defaults.baseURL;
    const response = await axios.post<AvatarResponse>(`${baseURL}/users/avatar`, formData, {
      withCredentials: true,
    });
    return response.data;
  },

  /**
   * アバター画像を削除し、デフォルトアバターに戻す。
   *
   * @throws {AxiosError} 401: 未認証
   */
  async deleteAvatar(): Promise<void> {
    await apiClient.delete("/users/avatar");
  },

  /**
   * サブスクリプション状態を取得する。
   *
   * @returns プラン名・料金・ステータス
   * @throws {AxiosError} 401: 未認証
   */
  async getSubscription(): Promise<SubscriptionResponse> {
    const response = await apiClient.get<SubscriptionResponse>("/subscriptions/me");
    return response.data;
  },

  /**
   * カスタマーポータルURLを生成する。
   *
   * @returns ポータルURL
   * @throws {AxiosError} 401: 未認証, 404: サブスクリプション未登録
   */
  async createPortal(): Promise<PortalResponse> {
    const response = await apiClient.post<PortalResponse>("/subscriptions/portal");
    return response.data;
  },

  /**
   * ニュースレター購読状態を取得する。
   *
   * @returns 購読情報（メールアドレス・is_active）
   * @throws {AxiosError} 401: 未認証, 404: 購読未登録
   */
  async getNewsletterSubscription(): Promise<NewsletterSubscriptionResponse> {
    const response = await apiClient.get<NewsletterSubscriptionResponse>(
      "/newsletter/subscription"
    );
    return response.data;
  },

  /**
   * ニュースレターを購読登録する。
   *
   * @param data - 購読登録リクエスト（email）
   * @returns 購読情報
   * @throws {AxiosError} 400: バリデーションエラー, 401: 未認証
   */
  async subscribeNewsletter(
    data: NewsletterSubscribeRequest
  ): Promise<NewsletterSubscriptionResponse> {
    const response = await apiClient.post<NewsletterSubscriptionResponse>(
      "/newsletter/subscribe",
      data
    );
    return response.data;
  },

  /**
   * ニュースレターの購読を解除する。
   *
   * @returns 解除後の購読情報（is_active=false）
   * @throws {AxiosError} 401: 未認証, 404: 購読未登録
   */
  async unsubscribeNewsletter(): Promise<NewsletterSubscriptionResponse> {
    const response =
      await apiClient.post<NewsletterSubscriptionResponse>("/newsletter/unsubscribe");
    return response.data;
  },

  /**
   * ニュースレター用メールアドレスを変更する。
   *
   * @param data - 変更リクエスト（email）
   * @returns 更新後の購読情報
   * @throws {AxiosError} 400: バリデーションエラー, 401: 未認証, 404: 購読未登録
   */
  async updateNewsletterEmail(
    data: NewsletterUpdateRequest
  ): Promise<NewsletterSubscriptionResponse> {
    const response = await apiClient.put<NewsletterSubscriptionResponse>(
      "/newsletter/subscription",
      data
    );
    return response.data;
  },
};
