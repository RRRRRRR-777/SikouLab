import { defineConfig } from "vitest/config";
import path from "path";
import react from "@vitejs/plugin-react";

/**
 * Vitest設定
 *
 * @description
 * - globals: true でグローバルテスト関数（describe, test, expect等）を有効化
 * - environment: jsdom でブラウザ環境をシミュレート
 * - include: テストファイルのパターンを指定
 * - exclude: 除外対象のディレクトリを指定
 * - alias: パスエイリアスを設定（@/ → ./）
 * - setupFiles: テスト実行前のセットアップファイル
 *
 * @see {@link https://vitest.dev/config/} 詳細設定
 * @example
 * ```bash
 * npm test           # 単体テスト実行
 * npm run test:watch # ウォッチモード（開発中）
 * ```
 */
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./"),
    },
  },
  test: {
    globals: true,
    environment: "jsdom",
    include: ["**/__tests__/**/*.{test,spec}.{js,ts,tsx}", "**/*.{test,spec}.{js,ts,tsx}"],
    exclude: ["node_modules", ".next", "dist"],
    setupFiles: ["./vitest.setup.ts"],
  },
});
