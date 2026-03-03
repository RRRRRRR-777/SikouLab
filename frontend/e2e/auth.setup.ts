import { test as setup, expect } from "@playwright/test";

const authFile = ".auth/user.json";

/**
 * E2E認証セットアップ（Firebase Auth Emulator）
 *
 * 1. Emulator REST APIでテストユーザーを作成
 * 2. アプリページに遷移（FirebaseはEmulatorに接続 + localStorage永続化済み）
 * 3. アプリのAuth instanceで直接サインイン（window.__E2E_SIGN_IN__）
 * 4. storageState保存（localStorage経由でセッション引き継ぎ）
 *
 * 前提条件:
 *   - Firebase Auth Emulator が localhost:9099 で起動済み
 *   - webServerが NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 付きで起動済み
 */
setup("authenticate", async ({ page }) => {
  const emulatorHost = "localhost:9099";
  const testEmail = "e2e-test@sicoulab.com";
  const testPassword = "e2e-test-password-123";

  // Emulator REST APIでテストユーザーを作成（既存ならスキップ）
  const signUpRes = await fetch(
    `http://${emulatorHost}/identitytoolkit.googleapis.com/v1/accounts:signUp?key=fake-api-key`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        email: testEmail,
        password: testPassword,
        returnSecureToken: true,
      }),
    },
  );

  if (!signUpRes.ok) {
    // ユーザーが既に存在する場合はサインインで確認
    const signInRes = await fetch(
      `http://${emulatorHost}/identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=fake-api-key`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: testEmail,
          password: testPassword,
          returnSecureToken: true,
        }),
      },
    );
    expect(signInRes.ok).toBeTruthy();
  }

  // アプリを読み込む（Firebase初期化 → Emulator接続 → localStorage永続化 → __E2E_SIGN_IN__公開）
  await page.goto("/login");
  await page.waitForLoadState("networkidle");

  // アプリのFirebaseが初期化されてヘルパーが公開されるまで待機
  await page.waitForFunction(
    () => typeof (window as unknown as Record<string, unknown>).__E2E_SIGN_IN__ === "function",
    null,
    { timeout: 15000 },
  );

  // onAuthStateChanged → authApi.login のフロー完了を待機するため、先にリスナーを登録
  const loginResponsePromise = page.waitForResponse(
    (response) => response.url().includes("/api/v1/auth/login") && response.status() === 200,
    { timeout: 15000 },
  );

  // アプリのAuth instance経由でサインイン（localStorage永続化）
  await page.evaluate(
    async ({ email, password }) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      await (window as any).__E2E_SIGN_IN__(email, password);
    },
    { email: testEmail, password: testPassword },
  );

  // 認証フロー完了を待機
  await loginResponsePromise;

  // セッション情報を保存（cookies + localStorage。Firebase auth tokenはlocalStorageに格納済み）
  await page.context().storageState({ path: authFile });
});
