"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutGrid, BookOpen, Newspaper, Vote, Settings, ChevronDown, X } from "lucide-react";
import { cn } from "@/lib/utils";

/** サイドバーのナビゲーション項目の型定義。 */
interface NavItem {
  label: string;
  href: string;
  icon: React.ReactNode;
  children?: { label: string; href: string }[];
}

const navItems: NavItem[] = [
  {
    label: "ダッシュボード",
    href: "/",
    icon: <LayoutGrid className="h-5 w-5" />,
  },
  {
    label: "記事",
    href: "/articles",
    icon: <BookOpen className="h-5 w-5" />,
    children: [{ label: "すべての記事", href: "/articles" }],
  },
  {
    label: "ニュース",
    href: "/news",
    icon: <Newspaper className="h-5 w-5" />,
  },
  {
    label: "投票",
    href: "/poll",
    icon: <Vote className="h-5 w-5" />,
  },
];

/**
 * アプリケーションのサイドバーコンポーネント。
 * デスクトップでは固定表示、モバイル/タブレットではスライドオーバーで表示する。
 *
 * @param open - コンポーネントのプロパティ
 * @param open.open - モバイル時のサイドバー開閉状態
 * @param open.setOpen - モバイル時のサイドバー開閉状態を変更するコールバック
 * @returns サイドバーコンポーネント
 */
export function Sidebar({ open, setOpen }: { open: boolean; setOpen: (open: boolean) => void }) {
  const pathname = usePathname();
  const [expandedItems, setExpandedItems] = useState<Record<string, boolean>>({});

  /**
   * 折りたたみ可能なナビ項目のトグル処理。
   *
   * @param label - トグルするナビ項目のラベル
   */
  const toggleExpand = (label: string) => {
    setExpandedItems((prev) => ({ ...prev, [label]: !prev[label] }));
  };

  const sidebarContent = (
    <div className="flex h-full flex-col">
      {/* ブランドロゴ */}
      <div className="flex h-16 items-center border-b px-6">
        <Link href="/" className="text-xl font-bold text-primary" onClick={() => setOpen(false)}>
          シコウラボ
        </Link>
        {/* モバイル時の閉じるボタン */}
        <button
          onClick={() => setOpen(false)}
          className="ml-auto lg:hidden"
          aria-label="サイドバーを閉じる"
        >
          <X className="h-5 w-5" />
        </button>
      </div>

      {/* メインナビゲーション */}
      <nav className="flex-1 space-y-1 px-3 py-4">
        {navItems.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.children?.some((child) => pathname === child.href) ?? false);
          const isExpanded = expandedItems[item.label] ?? false;

          return (
            <div key={item.label}>
              {item.children ? (
                // 折りたたみ対応のナビ項目
                <button
                  onClick={() => toggleExpand(item.label)}
                  className={cn(
                    "flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-primary/10 text-primary"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  {item.icon}
                  <span>{item.label}</span>
                  <ChevronDown
                    className={cn(
                      "ml-auto h-4 w-4 transition-transform",
                      isExpanded && "rotate-180"
                    )}
                  />
                </button>
              ) : (
                // 通常のナビ項目
                <Link
                  href={item.href}
                  onClick={() => setOpen(false)}
                  className={cn(
                    "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-primary/10 text-primary"
                      : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  {item.icon}
                  <span>{item.label}</span>
                </Link>
              )}

              {/* 子メニュー */}
              {item.children && isExpanded && (
                <div className="ml-8 mt-1 space-y-1">
                  {item.children.map((child) => (
                    <Link
                      key={child.href}
                      href={child.href}
                      onClick={() => setOpen(false)}
                      className={cn(
                        "block rounded-md px-3 py-1.5 text-sm transition-colors",
                        pathname === child.href
                          ? "text-primary font-medium"
                          : "text-muted-foreground hover:text-accent-foreground"
                      )}
                    >
                      {child.label}
                    </Link>
                  ))}
                </div>
              )}
            </div>
          );
        })}
      </nav>

      {/* ユーザーセクション */}
      <div className="border-t p-4">
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded-full bg-muted" />
          <div className="flex-1 text-sm">
            <p className="font-medium">ユーザー名</p>
          </div>
          <Link
            href="/settings"
            onClick={() => setOpen(false)}
            className="text-muted-foreground hover:text-foreground"
            aria-label="設定"
          >
            <Settings className="h-5 w-5" />
          </Link>
        </div>
      </div>
    </div>
  );

  return (
    <>
      {/* デスクトップ: 固定サイドバー */}
      <aside className="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-[280px] lg:flex-col border-r bg-background">
        {sidebarContent}
      </aside>

      {/* モバイル/タブレット: オーバーレイ + スライドオーバー */}
      {open && (
        <>
          <div
            className="fixed inset-0 z-40 bg-black/50 lg:hidden"
            onClick={() => setOpen(false)}
          />
          <aside className="fixed inset-y-0 left-0 z-50 w-[280px] bg-background lg:hidden">
            {sidebarContent}
          </aside>
        </>
      )}
    </>
  );
}
