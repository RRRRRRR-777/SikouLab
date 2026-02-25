import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  // Turbopackのルートディレクトリを明示的に指定
  turbopack: {
    root: process.cwd(),
  },
};

export default nextConfig;
