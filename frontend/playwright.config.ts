import { defineConfig, devices } from "@playwright/test";

/**
 * E2Eテスト用ポート
 *
 * 開発サーバー（:3000）と競合しないよう別ポートで起動する。
 * webServerが NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST 付きで起動するため、
 * Emulator接続済みのフロントエンドに対してテストが実行される。
 */
const E2E_PORT = 3100;
const baseURL = process.env.E2E_BASE_URL ?? `http://localhost:${E2E_PORT}`;

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",

  use: {
    baseURL,
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },

  // E2E用にEmulator接続済みのNext.jsを自動起動
  webServer: {
    command: `NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 NEXT_DIST_DIR=.next-e2e npm run dev -- --port ${E2E_PORT}`,
    port: E2E_PORT,
    reuseExistingServer: !process.env.CI,
    timeout: 30_000,
  },

  projects: [
    // 認証セットアップ（Emulator経由でログイン → storageState保存）
    {
      name: "setup",
      testMatch: /auth\.setup\.ts/,
    },

    // モバイル（モバイルファースト: 最優先）
    {
      name: "mobile",
      use: {
        ...devices["iPhone 14"],
        storageState: ".auth/user.json",
      },
      dependencies: ["setup"],
    },

    // デスクトップ
    {
      name: "desktop",
      use: {
        ...devices["Desktop Chrome"],
        storageState: ".auth/user.json",
      },
      dependencies: ["setup"],
    },
  ],
});
