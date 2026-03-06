/**
 * 設定画面メインコンテナ
 *
 * @description
 * 設定画面の4セクション（プロフィール、サブスクリプション、メール、FAQ）を統合する。
 * MainLayoutの中で表示され、各セクションは独立してデータ取得を行う。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 詳細設計書
 */

"use client";

import { ProfileSection } from "./ProfileSection";
import { SubscriptionSection } from "./SubscriptionSection";
import { EmailSection } from "./EmailSection";
import { FaqSection } from "./FaqSection";

/**
 * 設定画面メインコンテナコンポーネント
 *
 * 4つのセクションを縦に並べた設定画面を表示する。
 *
 * @returns 設定画面
 */
export function SettingsPage() {
  return (
    <div className="mx-auto max-w-3xl">
      <h1 className="text-2xl font-bold text-[var(--color-text)]">設定</h1>

      <div className="mt-6 space-y-6">
        <ProfileSection />
        <SubscriptionSection />
        <EmailSection />
        <FaqSection />
      </div>
    </div>
  );
}
