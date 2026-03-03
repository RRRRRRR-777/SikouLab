/**
 * Firebase Authentication 初期化・設定
 *
 * @description
 * Firebase JS SDKの初期化と認証プロバイダ設定を行う。
 *
 * @module lib/auth/firebase
 */

import { initializeApp, getApps, FirebaseApp } from "firebase/app";
import {
  getAuth,
  connectAuthEmulator,
  setPersistence,
  browserLocalPersistence,
  signInWithEmailAndPassword,
  Auth,
  GoogleAuthProvider,
  OAuthProvider,
  signInWithPopup,
  signOut as firebaseSignOut,
  User as FirebaseUser,
  onAuthStateChanged,
} from "firebase/auth";

/**
 * Firebase Config
 *
 * 環境変数から設定を読み込む。
 * Firebaseコンソール → プロジェクトの設定 → 全般 → SDKの設定と構成から取得。
 * signInWithPopup などのOAuth認証には authDomain が必須。
 */
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY ?? "",
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN ?? "",
};

/**
 * Firebaseアプリインスタンス（シングルトン）
 */
let firebaseApp: FirebaseApp | null = null;

/**
 * Firebase Authインスタンス（シングルトン）
 */
let authInstance: Auth | null = null;

/**
 * Firebaseアプリを初期化する
 *
 * @returns Firebaseアプリインスタンス
 */
export function initializeFirebaseApp(): FirebaseApp {
  if (firebaseApp) {
    return firebaseApp;
  }

  const apps = getApps();
  if (apps.length > 0) {
    firebaseApp = apps[0];
  } else {
    if (!firebaseConfig.apiKey) {
      throw new Error("Firebase API Key is not configured. Please set NEXT_PUBLIC_FIREBASE_API_KEY.");
    }
    firebaseApp = initializeApp(firebaseConfig);
  }

  return firebaseApp;
}

/**
 * Firebase Auth EmulatorのURL
 *
 * E2Eテスト時にFirebase Auth Emulatorを使用するための設定。
 * NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST が設定されている場合にEmulatorに接続する。
 */
const EMULATOR_HOST = process.env.NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST ?? "";

/**
 * Firebase Authインスタンスを取得する
 *
 * NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST が設定されている場合、
 * Firebase Auth Emulatorに接続する（E2Eテスト用）。
 *
 * @returns Firebase Authインスタンス
 */
export function getFirebaseAuth(): Auth {
  if (authInstance) {
    return authInstance;
  }

  initializeFirebaseApp();
  authInstance = getAuth();

  // E2Eテスト用: Firebase Auth Emulatorに接続（本番環境では無効化）
  if (EMULATOR_HOST && process.env.NODE_ENV !== "production") {
    connectAuthEmulator(authInstance, `http://${EMULATOR_HOST}`, { disableWarnings: true });

    // storageStateでセッションを引き継ぐため、IndexedDBではなくlocalStorageに永続化
    const persistenceReady = setPersistence(authInstance, browserLocalPersistence);

    // ブラウザ上でE2Eテストからサインインできるヘルパーを公開
    if (typeof window !== "undefined") {
      (window as unknown as Record<string, unknown>).__E2E_SIGN_IN__ = async (
        email: string,
        password: string,
      ) => {
        await persistenceReady;
        return signInWithEmailAndPassword(authInstance!, email, password);
      };
    }
  }

  return authInstance;
}

/**
 * Googleプロバイダでログインする
 *
 * @returns Firebase認証結果
 */
export async function signInWithGoogle(): Promise<FirebaseUser> {
  const auth = getFirebaseAuth();
  const provider = new GoogleAuthProvider();
  const result = await signInWithPopup(auth, provider);
  return result.user;
}

/**
 * Appleプロバイダでログインする
 *
 * @returns Firebase認証結果
 */
export async function signInWithApple(): Promise<FirebaseUser> {
  const auth = getFirebaseAuth();
  const provider = new OAuthProvider("apple.com");
  const result = await signInWithPopup(auth, provider);
  return result.user;
}

/**
 * X（Twitter）プロバイダでログインする
 *
 * @returns Firebase認証結果
 */
export async function signInWithX(): Promise<FirebaseUser> {
  const auth = getFirebaseAuth();
  const provider = new OAuthProvider("twitter.com");
  const result = await signInWithPopup(auth, provider);
  return result.user;
}

/**
 * ログアウトする
 *
 * @returns ログアウト結果
 */
export async function signOut(): Promise<void> {
  const auth = getFirebaseAuth();
  await firebaseSignOut(auth);
}

/**
 * 認証状態の変化を監視する
 *
 * @param callback - 認証状態変化時のコールバック
 * @returns 監視解除関数
 */
export function onAuthStateChangedHelper(
  callback: (user: FirebaseUser | null) => void,
): () => void {
  const auth = getFirebaseAuth();
  return onAuthStateChanged(auth, callback);
}

/**
 * FirebaseユーザーからID Tokenを取得する
 *
 * @param user - Firebaseユーザー
 * @param forceRefresh - trueの場合、キャッシュを無視して最新トークンを取得する
 * @returns ID Token
 */
export async function getIdToken(user: FirebaseUser, forceRefresh = false): Promise<string> {
  return await user.getIdToken(forceRefresh);
}
