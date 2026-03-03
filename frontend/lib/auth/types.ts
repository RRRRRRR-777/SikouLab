/**
 * 認証関連の型定義
 *
 * @description
 * Firebase認証とバックエンドAPIで使用する型を定義する。
 *
 * @module lib/auth/types
 */

/**
 * バックエンドAPIから返されるユーザー情報
 */
export interface User {
  /** ユーザーID */
  id: number;
  /** OAuthプロバイダ識別子（例: "google.com", "apple.com"） */
  oauth_provider: string;
  /** 表示名 */
  name: string;
  /** OAuthプロバイダから取得した表示名 */
  display_name: string | null;
  /** アバター画像URL */
  avatar_url: string | null;
  /** ユーザーロール */
  role: "admin" | "writer" | "user";
  /** プランID */
  plan_id: number | null;
  /** サブスクリプション状態 */
  subscription_status: "active" | "canceled" | "past_due" | null;
  /** 作成日時 */
  created_at: string;
  /** 更新日時 */
  updated_at: string;
}

/**
 * フロントエンドで使用するユーザー情報（キャメルケース）
 */
export interface AuthUser {
  /** ユーザーID */
  id: number;
  /** OAuthプロバイダ識別子 */
  oauthProvider: string;
  /** 表示名 */
  name: string;
  /** OAuthプロバイダから取得した表示名 */
  displayName: string | null;
  /** アバター画像URL */
  avatarUrl: string | null;
  /** ユーザーロール */
  role: "admin" | "writer" | "user";
  /** プランID */
  planId: number | null;
  /** サブスクリプション状態 */
  subscriptionStatus: "active" | "canceled" | "past_due" | null;
  /** 作成日時 */
  createdAt: string;
  /** 更新日時 */
  updatedAt: string;
}

/**
 * ログインAPIレスポンス
 */
export interface LoginResponse {
  /** ユーザー情報 */
  user: User;
  /** 初回ログインかどうか（trueの場合はサブスクリプション登録へ誘導） */
  is_first_login: boolean;
}

/**
 * セッション確認APIレスポンス
 */
export interface MeResponse {
  /** ユーザー情報 */
  user: User;
}

/**
 * APIエラーレスポンス
 */
export interface ApiError {
  /** エラーコード */
  code: string;
  /** エラーメッセージ */
  message: string;
}

/**
 * OAuthプロバイダの種類
 */
export type OAuthProvider = "google" | "apple" | "twitter.com";

/**
 * スネークケースのUserをキャメルケースのAuthUserに変換する
 *
 * @param user - バックエンドAPIから返されるユーザー情報
 * @returns フロントエンドで使用するユーザー情報
 */
export function toAuthUser(user: User): AuthUser {
  return {
    id: user.id,
    oauthProvider: user.oauth_provider,
    name: user.name,
    displayName: user.display_name,
    avatarUrl: user.avatar_url,
    role: user.role,
    planId: user.plan_id,
    subscriptionStatus: user.subscription_status,
    createdAt: user.created_at,
    updatedAt: user.updated_at,
  };
}
