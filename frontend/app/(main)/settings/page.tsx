/**
 * 設定画面ページ（/settings）
 *
 * @description
 * MainLayout内で設定画面を表示する。
 * Client Componentの SettingsPage をインポートして描画する。
 */

import { SettingsPage } from "@/components/settings/SettingsPage";

/**
 * 設定ページのメタデータ
 */
export const metadata = {
  title: "設定 | シコウラボ",
};

/**
 * 設定画面ページコンポーネント
 *
 * @returns 設定画面
 */
export default function SettingsPageRoute() {
  return <SettingsPage />;
}
