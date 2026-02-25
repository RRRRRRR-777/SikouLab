"use client";

import { useState } from "react";
import { Menu, Search, X } from "lucide-react";

/**
 * アプリケーションのヘッダーコンポーネント。
 * モバイルではハンバーガーメニューとタイトルを表示し、検索は右下のフローティングボタンで提供する。
 * デスクトップではヘッダー内に検索バーを表示する。
 *
 * @param setOpen - サイドバーの開閉状態を変更するコールバック
 */
export function Header({ setOpen }: { setOpen: (open: boolean) => void }) {
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");

  /**
   * ×ボタンの動作を統合する。
   * テキストあり → 入力をクリア、テキストなし → 検索バーを閉じる
   */
  const handleXClick = () => {
    if (searchValue) {
      setSearchValue("");
    } else {
      setSearchOpen(false);
    }
  };

  return (
    <>
      <header className="sticky top-0 z-30 flex h-16 items-center gap-4 border-b bg-background px-4 lg:px-6">
        {/* ハンバーガーメニュー。デスクトップではサイドナビが常時表示のため不要 */}
        <button onClick={() => setOpen(true)} className="lg:hidden" aria-label="メニューを開く">
          <Menu className="h-6 w-6" />
        </button>

        <span className="text-lg font-bold text-primary lg:hidden">シコウラボ</span>

        {/* デスクトップ: 検索バー。モバイルでは右下フローティングボタンで代替 */}
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

      {/* モバイル: 検索UI。デスクトップではヘッダー内検索バーで代替 */}
      <div className="lg:hidden">
        {searchOpen ? (
          /* 展開時: 左右均等余白で画面幅いっぱいに広げる */
          <div className="fixed bottom-6 left-4 right-4 z-40 flex items-center gap-3 rounded-full border bg-background px-4 py-3 shadow-lg">
            <Search className="h-4 w-4 shrink-0 text-muted-foreground" />
            <input
              type="text"
              value={searchValue}
              onChange={(e) => setSearchValue(e.target.value)}
              placeholder="検索..."
              autoFocus
              className="flex-1 bg-transparent text-sm outline-none"
            />
            {/* テキストあり→クリア、テキストなし→閉じる の1ボタン統合 */}
            <button
              onClick={handleXClick}
              aria-label={searchValue ? "入力をクリア" : "検索を閉じる"}
            >
              <X className="h-4 w-4 text-muted-foreground" />
            </button>
          </div>
        ) : (
          <button
            onClick={() => setSearchOpen(true)}
            aria-label="検索を開く"
            className="fixed bottom-6 right-6 z-40 flex h-14 w-14 items-center justify-center rounded-full bg-[#E86D00] text-white shadow-lg"
          >
            <Search className="h-6 w-6" />
          </button>
        )}
      </div>
    </>
  );
}
