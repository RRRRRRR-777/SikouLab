/**
 * サンプルテスト - Vitest動作確認用
 */
import { describe, it, expect } from "vitest";

describe("サンプルテスト", () => {
  it("足し算が正しく動作すること", () => {
    expect(1 + 1).toBe(2);
  });

  it("文字列結合が正しく動作すること", () => {
    expect("Hello" + " " + "World").toBe("Hello World");
  });
});
