import { defineConfig } from "vitest/config";

/**
 * Vitest設定
 * - globals: true でグローバルテスト関数（describe, test, expect等）を有効化
 * - environment: jsdom でブラウザ環境をシミュレート
 */
export default defineConfig({
  test: {
    globals: true,
    environment: "jsdom",
    include: ["**/__tests__/**/*.{test,spec}.{js,ts,tsx}", "**/*.{test,spec}.{js,ts,tsx}"],
    exclude: ["node_modules", ".next", "dist"],
  },
});
