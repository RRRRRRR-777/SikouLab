import { test as setup, expect } from "@playwright/test";

const authFile = ".auth/user.json";

/**
 * E2E認証セットアップ（Firebase Auth Emulator）
 *
 * 1. Emulator REST APIでテストユーザーを作成
 * 2. アプリページに遷移（FirebaseはEmulatorに接続 + localStorage永続化済み）
 * 3. アプリのAuth instanceで直接サインイン（window.__E2E_SIGN_IN__）
 * 4. 取得したID Tokenで直接バックエンドログインAPIを呼び、セッションCookieを取得
 * 5. storageState保存（cookies + localStorage経由でセッション引き継ぎ）
 *
 * 注意: onAuthStateChangedは/loginページではauthApi.loginをスキップするため、
 * E2Eではpage.request.postで直接バックエンドを呼ぶ。
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

  // Firebase Auth経由でサインインし、ID Tokenを取得する。
  // onAuthStateChangedは/loginページではauthApi.loginをスキップするため、
  // ID Tokenを返してもらい、直接バックエンドのログインAPIを呼ぶ。
  const idToken = await page.evaluate(
    async ({ email, password }) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const { user } = await (window as any).__E2E_SIGN_IN__(email, password);
      return await user.getIdToken();
    },
    { email: testEmail, password: testPassword },
  );

  // バックエンドのログインAPIを直接呼び、セッションCookieを取得する
  const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";
  const loginRes = await page.request.post(`${apiUrl}/auth/login`, {
    data: { id_token: idToken },
  });
  expect(loginRes.ok()).toBeTruthy();

  // セッション情報を保存（cookies + localStorage。Firebase auth tokenはlocalStorageに格納済み）
  await page.context().storageState({ path: authFile });
});
