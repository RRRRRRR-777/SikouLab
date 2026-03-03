/**
 * 認証APIクライアント
 *
 * @description
 * バックエンド認証APIと通信するクライアント。
 * Firebase ID Tokenを使用したログイン、セッション確認、ログアウトを行う。
 *
 * @module lib/auth/auth-api
 */

import { apiClient } from "@/lib/api";
import type { LoginResponse, MeResponse } from "./types";

/**
 * 認証API操作
 */
export const authApi = {
  /**
   * ログイン（Firebase ID Tokenを検証）
   *
   * @param idToken - Firebase ID Token
   * @returns ログインレスポンス（ユーザー情報、初回ログインフラグ）
   * @throws {AxiosError} 認証失敗時
   */
  async login(idToken: string): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>("/auth/login", {
      id_token: idToken,
    });
    return response.data;
  },

  /**
   * セッション確認・ユーザー情報取得
   *
   * @returns ユーザー情報
   * @throws {AxiosError} 未認証時
   */
  async getMe(): Promise<MeResponse> {
    const response = await apiClient.get<MeResponse>("/auth/me");
    return response.data;
  },

  /**
   * ログアウト（セッションCookieを削除）
   *
   * @throws {AxiosError} ログアウト失敗時
   */
  async logout(): Promise<void> {
    await apiClient.post("/auth/logout");
  },
};
