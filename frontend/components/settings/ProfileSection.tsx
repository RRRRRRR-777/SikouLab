/**
 * プロフィール設定セクション（F-10-1）
 *
 * @description
 * アバター表示/変更/削除、ユーザーID表示（読み取り専用）、表示名編集を提供する。
 * auth contextからユーザー情報を取得し、settings APIで更新する。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 詳細設計書
 */

"use client";

import { useState, useRef } from "react";
import { toast } from "sonner";
import { useAuth } from "@/lib/auth/auth-context";
import { settingsApi } from "@/lib/settings/settings-api";

/**
 * プロフィール設定セクションコンポーネント
 *
 * アバター管理と表示名編集のUIを提供する。
 *
 * @returns プロフィール設定セクション
 */
export function ProfileSection() {
  const { user, isLoading, refresh } = useAuth();
  const [displayName, setDisplayName] = useState(user?.displayName ?? user?.name ?? "");
  const [isSaving, setIsSaving] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  /**
   * 表示名を保存する。
   */
  const handleSaveDisplayName = async () => {
    if (!displayName.trim()) {
      toast.error("表示名を入力してください");
      return;
    }

    setIsSaving(true);
    try {
      await settingsApi.updateProfile({ display_name: displayName });
      await refresh();
      toast.success("保存しました");
    } catch {
      toast.error("保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  /**
   * アバター画像を変更する。ファイル選択後に自動アップロードする。
   *
   * @param e - ファイル選択イベント
   */
  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      await settingsApi.uploadAvatar(file);
      await refresh();
      toast.success("アバターを変更しました");
    } catch {
      toast.error("アバターの変更に失敗しました");
    }
  };

  /**
   * アバター画像を削除してデフォルトに戻す。
   */
  const handleDeleteAvatar = async () => {
    try {
      await settingsApi.deleteAvatar();
      await refresh();
      toast.success("アバターを削除しました");
    } catch {
      toast.error("アバターの削除に失敗しました");
    }
  };

  if (isLoading) {
    return (
      <section className="rounded-lg border border-gray-200 p-4 dark:border-gray-800 lg:p-6">
        <h2 className="text-xl font-bold text-[var(--color-text)]">プロフィール設定</h2>
        <div className="mt-4 flex flex-col gap-4 md:flex-row md:items-start md:gap-6">
          <div className="flex flex-col items-center gap-2">
            <div className="h-20 w-20 animate-pulse rounded-full bg-gray-200 dark:bg-gray-700" />
          </div>
          <div className="flex flex-1 flex-col gap-4">
            <div className="h-6 w-32 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
            <div className="h-10 w-full animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="rounded-lg border border-gray-200 p-4 dark:border-gray-800 lg:p-6">
      <h2 className="text-xl font-bold text-[var(--color-text)]">プロフィール設定</h2>

      <div className="mt-4 flex flex-col gap-4 md:flex-row md:items-end md:gap-6">
        {/* アバター */}
        <div className="flex flex-col items-center gap-2 md:shrink-0">
          {user?.avatarUrl ? (
            /* eslint-disable-next-line @next/next/no-img-element */
            <img
              src={user.avatarUrl}
              alt="アバター"
              className="h-20 w-20 rounded-full object-cover"
            />
          ) : (
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-gray-200 dark:bg-gray-700">
              <span className="text-2xl text-[var(--color-muted-foreground)]">
                {(user?.displayName ?? user?.name ?? "?").charAt(0)}
              </span>
            </div>
          )}
          <div className="flex gap-2">
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="min-h-[44px] rounded-md bg-[var(--color-primary)] px-3 py-2 text-lg text-white hover:opacity-80"
            >
              変更
            </button>
            {user?.avatarUrl && (
              <button
                type="button"
                onClick={handleDeleteAvatar}
                className="min-h-[44px] rounded-md border border-gray-300 px-3 py-2 text-lg text-[var(--color-text)] hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-800"
              >
                削除
              </button>
            )}
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png,image/gif"
            onChange={handleAvatarChange}
            className="hidden"
          />
        </div>

        {/* ユーザー情報 */}
        <div className="flex flex-1 flex-col gap-4">
          {/* ユーザーID（読み取り専用） */}
          <div>
            <label className="text-lg font-bold text-[var(--color-text)]">ユーザーID</label>
            <p className="mt-1 text-lg font-medium text-[var(--color-text-secondary)]">{user?.oauthUserId ?? ""}</p>
          </div>

          {/* 表示名 */}
          <div>
            <label
              htmlFor="display-name"
              className="text-lg font-bold text-[var(--color-text)]"
            >
              表示名
            </label>
            <div className="mt-1 flex gap-2">
              <input
                id="display-name"
                type="text"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                maxLength={50}
                className="min-h-[44px] flex-1 rounded-md border border-gray-300 bg-[var(--color-bg)] px-3 py-2 text-lg text-[var(--color-text)] outline-none focus:ring-2 focus:ring-[var(--color-primary)] dark:border-gray-600"
              />
              <button
                type="button"
                onClick={handleSaveDisplayName}
                disabled={isSaving}
                className="min-h-[44px] rounded-md bg-[var(--color-primary)] px-4 py-2 text-lg text-white hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-50"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
