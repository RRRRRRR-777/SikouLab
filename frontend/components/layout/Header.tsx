"use client";

import { Menu, Search } from "lucide-react";

/**
 * アプリケーションのヘッダーコンポーネント。
 * モバイルではハンバーガーメニューとタイトル、デスクトップでは検索バーを表示する。
 *
 * @param setOpen - サイドバーの開閉状態を変更するコールバック
 */
export function Header({ setOpen }: { setOpen: (open: boolean) => void }) {
  return (
    <header className="sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-background px-4 lg:px-6">
      {/* モバイル: ハンバーガーメニュー */}
      <button onClick={() => setOpen(true)} className="lg:hidden" aria-label="メニューを開く">
        <Menu className="h-6 w-6" />
      </button>

      {/* モバイル: タイトル */}
      <span className="text-lg font-bold text-primary lg:hidden">SicouLab</span>

      {/* デスクトップ: 検索バー */}
      <div className="hidden lg:flex lg:flex-1">
        <div className="relative w-full max-w-md">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            type="search"
            placeholder="検索..."
            className="h-10 w-full rounded-md border border-input bg-background pl-10 pr-4 text-sm outline-none focus:ring-2 focus:ring-ring"
          />
        </div>
      </div>
    </header>
  );
}
