"use client";

import { useState } from "react";
import { Sidebar } from "@/components/layout/Sidebar";
import { Header } from "@/components/layout/Header";

/**
 * 認証が必要なページの共通レイアウト。
 * サイドバーとヘッダーを含むメインレイアウトを提供する。
 */
export default function MainLayout({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="min-h-screen">
      <Sidebar open={sidebarOpen} setOpen={setSidebarOpen} />
      <div className="lg:pl-[280px]">
        <Header setOpen={setSidebarOpen} />
        <main className="p-4 lg:p-6">{children}</main>
      </div>
    </div>
  );
}
