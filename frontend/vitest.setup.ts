/**
 * Vitestセットアップファイル
 *
 * @description
 * - @testing-library/jest-domのマッチャーをグローバルに追加
 * - jsdom環境のデフォルト設定
 */

import { expect, afterEach } from "vitest";
import { cleanup } from "@testing-library/react";
import * as matchers from "@testing-library/jest-dom/matchers";
import "@testing-library/jest-dom/vitest";

// 各テスト後にDOMをクリーンアップ
afterEach(() => {
  cleanup();
});

// jest-domのマッチャーを拡張
expect.extend(matchers);
