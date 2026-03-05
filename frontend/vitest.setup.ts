/**
 * Vitestセットアップファイル
 *
 * @description
 * - @testing-library/jest-domのマッチャーをグローバルに追加
 * - jsdom環境のデフォルト設定
 * - next-themes の ThemeProvider が使用する window.matchMedia のスタブを定義
 */

import { expect, afterEach } from "vitest";
import { cleanup } from "@testing-library/react";
import * as matchers from "@testing-library/jest-dom/matchers";
import "@testing-library/jest-dom/vitest";

// jsdom は window.matchMedia を実装していないため、next-themes 向けにスタブを定義する
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: (query: string): MediaQueryList => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});

// 各テスト後にDOMをクリーンアップ
afterEach(() => {
  cleanup();
});

// jest-domのマッチャーを拡張
expect.extend(matchers);
