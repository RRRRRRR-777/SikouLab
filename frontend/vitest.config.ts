import { defineConfig } from "vitest/config";
// import nextjs from '@next/eslint-plugin-next'; // 未使用のためコメント化

export default defineConfig({
  test: {
    globals: true,
    environment: "jsdom",
    include: ["**/__tests__/**/*.{test,spec}.{js,ts,tsx}", "**/*.{test,spec}.{js,ts,tsx}"],
    exclude: ["node_modules", ".next", "dist"],
  },
});
