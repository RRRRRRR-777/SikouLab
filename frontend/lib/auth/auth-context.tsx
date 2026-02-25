/**
 * 認証コンテキスト
 *
 * @description
 * 認証状態を管理し、認証操作を提供する。
 * Google / Apple / X によるOAuth認証をサポートする。
 *
 * @module lib/auth/auth-context
 */

"use client";

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { useRouter, usePathname } from "next/navigation";
import { toast } from "sonner";
import {
  signInWithGoogle as firebaseSignInGoogle,
  signInWithApple as firebaseSignInApple,
  signInWithX as firebaseSignInX,
  signOut as firebaseSignOut,
  getIdToken,
  onAuthStateChangedHelper,
} from "./firebase";
import type { User as FirebaseUser } from "firebase/auth";
import { authApi } from "./auth-api";
import { toAuthUser, type AuthUser } from "./types";
import { logError, logWarn } from "../logger";

/**
 * 認証コンテキストの型定義
 */
interface AuthContextValue {
  /** 認証済みユーザー情報 */
  user: AuthUser | null;
  /** 認証済みかどうか */
  isAuthenticated: boolean;
  /** ローディング中かどうか */
  isLoading: boolean;
  /** Googleでログイン */
  loginWithGoogle: () => Promise<void>;
  /** Appleでログイン */
  loginWithApple: () => Promise<void>;
  /** Xでログイン */
  loginWithX: () => Promise<void>;
  /** ログアウト */
  logout: () => Promise<void>;
  /** ユーザー情報を再取得 */
  refresh: () => Promise<void>;
}

/**
 * 認証コンテキストのデフォルト値
 */
const defaultValue: AuthContextValue = {
  user: null,
  isAuthenticated: false,
  isLoading: true,
  loginWithGoogle: async () => {},
  loginWithApple: async () => {},
  loginWithX: async () => {},
  logout: async () => {},
  refresh: async () => {},
};

/**
 * 認証コンテキスト
 */
const AuthContext = createContext<AuthContextValue>(defaultValue);

/**
 * 認証プロバイダーのプロパティ
 */
interface AuthProviderProps {
  children: ReactNode;
}

/**
 * 認証プロバイダー
 *
 * Firebase認証状態を監視し、認証操作を提供する。
 *
 * @param root0 - コンポーネントのプロパティ
 * @param root0.children - 子コンポーネント
 * @returns 認証コンテキストを提供するプロバイダー
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const router = useRouter();
  const pathname = usePathname();
  // レンダーのたびに最新のパス名を参照できるよう更新する
  const pathnameRef = useRef(pathname);
  pathnameRef.current = pathname;
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  /**
   * ユーザー情報をバックエンドから取得する
   */
  const fetchUserInfo = async (): Promise<void> => {
    try {
      const response = await authApi.getMe();
      setUser(toAuthUser(response.user));
    } catch (error) {
      // 未認証の場合は警告ログを出力
      logWarn("Failed to fetch user info", { error });
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Firebase認証→バックエンド認証のフローを実行する
   *
   * @param signInFn - OAuthプロバイダのサインイン関数
   */
  const authenticate = async (
    signInFn: () => Promise<FirebaseUser>,
  ): Promise<void> => {
    setIsLoading(true);
    try {
      // Firebase認証
      const firebaseUser = await signInFn();

      // ID Tokenを取得
      const idToken = await getIdToken(firebaseUser);

      // バックエンドで認証
      const loginResponse = await authApi.login(idToken);
      const authUser = toAuthUser(loginResponse.user);

      setUser(authUser);

      // 初回ログイン時はサブスクリプション画面へ、既存ユーザーはダッシュボードへ
      if (loginResponse.is_first_login) {
        router.push("/subscription");
      } else {
        router.push("/");
      }

      toast.success("ログインしました");
    } catch (error) {
      // エラーの詳細情報を出力
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { data?: { code?: string; message?: string }; status?: number }; request?: unknown; message?: string };
        logError("Authentication failed", {
          status: axiosError.response?.status,
          data: axiosError.response?.data,
          message: axiosError.message,
        });
        // バックエンドからのエラーレスポンスがあれば詳細メッセージを表示
        if (axiosError.response?.data?.message) {
          toast.error(axiosError.response.data.message);
        } else {
          toast.error("ログインに失敗しました");
        }
      } else {
        logError("Authentication failed", { error });
        toast.error("ログインに失敗しました");
      }
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Googleでログインする
   */
  const loginWithGoogle = async (): Promise<void> => {
    await authenticate(() => firebaseSignInGoogle());
  };

  /**
   * Appleでログインする
   */
  const loginWithApple = async (): Promise<void> => {
    await authenticate(() => firebaseSignInApple());
  };

  /**
   * Xでログインする
   */
  const loginWithX = async (): Promise<void> => {
    await authenticate(() => firebaseSignInX());
  };

  /**
   * ログアウトする
   */
  const logout = async (): Promise<void> => {
    setIsLoading(true);
    try {
      // Firebaseサインアウト
      await firebaseSignOut();

      // バックエンドでログアウト
      await authApi.logout();

      setUser(null);
      router.push("/login");
      toast.success("ログアウトしました");
    } catch (error) {
      logError("Logout failed", { error });
      toast.error("ログアウトに失敗しました");
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * ユーザー情報を再取得する
   */
  const refresh = async (): Promise<void> => {
    await fetchUserInfo();
  };

  /**
   * Firebase認証状態の監視を開始する
   */
  useEffect(() => {
    const unsubscribe = onAuthStateChangedHelper(async (firebaseUser) => {
      if (firebaseUser && pathnameRef.current !== "/login") {
        // Firebase認証済みかつログインページ以外の場合、バックエンドでユーザー情報を取得
        await fetchUserInfo();
      } else {
        // 未認証状態またはログインページ（/auth/meを呼ぶ必要がない）
        setUser(null);
        setIsLoading(false);
      }
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const value: AuthContextValue = {
    user,
    isAuthenticated: user !== null,
    isLoading,
    loginWithGoogle,
    loginWithApple,
    loginWithX,
    logout,
    refresh,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * 認証コンテキストを使用するフック
 *
 * @returns 認証コンテキストの値
 * @example
 * ```tsx
 * const { user, isAuthenticated, loginWithGoogle } = useAuth();
 * ```
 */
export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
