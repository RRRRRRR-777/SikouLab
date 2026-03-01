import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  // E2Eテスト時は別ディレクトリを使用し、開発サーバーとの競合を回避
  ...(process.env.NEXT_DIST_DIR && { distDir: process.env.NEXT_DIST_DIR }),
  // Turbopackのルートディレクトリを明示的に指定
  turbopack: {
    root: process.cwd(),
  },
};

export default nextConfig;
