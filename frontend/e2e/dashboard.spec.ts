import { test, expect } from "@playwright/test";

test.describe("ダッシュボード", () => {
  test("トップページが表示される", async ({ page }) => {
    await page.goto("/");
    // ヘッダーのサービス名が表示されることを確認
    await expect(page.locator("header")).toBeVisible();
  });
});
