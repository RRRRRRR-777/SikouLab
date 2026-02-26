"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "next-themes";
import { useState } from "react";
import { AuthProvider } from "../lib/auth/auth-context";

/**
 * アプリケーション全体のProviderをまとめるコンポーネント。
 * TanStack Query、ダークモード（next-themes）、認証（AuthProvider）を提供する。
 * AuthProvider は QueryClientProvider の内側に配置し、認証フック内で TanStack Query が使えるようにする。
 */
export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
          },
        },
      })
  );

  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>{children}</AuthProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
}
