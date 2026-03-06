import type { NextConfig } from "next";

const isProduction = process.env.NODE_ENV === "production";

const nextConfig: NextConfig = {
  output: "standalone",
  // E2Eテスト時は別ディレクトリを使用し、開発サーバーとの競合を回避
  ...(process.env.NEXT_DIST_DIR && { distDir: process.env.NEXT_DIST_DIR }),
  // Turbopackのルートディレクトリを明示的に指定
  turbopack: {
    root: process.cwd(),
  },
  // 画像ドメインの許可（本番: GCSのみ、開発: GCSエミュレータも追加）
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "storage.googleapis.com",
      },
      // 開発環境ではGCSエミュレータ（localhost:4443）を許可
      ...(!isProduction
        ? [
            {
              protocol: "http" as const,
              hostname: "localhost",
              port: "4443",
            },
          ]
        : []),
    ],
  },
  // 開発環境のみ: GCSエミュレータへのリクエストをプロキシ（CORS/ORB回避）
  async rewrites() {
    if (isProduction) {
      return [];
    }
    return [
      {
        source: "/storage/:path*",
        destination: "http://localhost:4443/storage/:path*",
      },
    ];
  },
};

export default nextConfig;
