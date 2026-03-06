/**
 * FAQ・問い合わせセクション（F-10-4）
 *
 * @description
 * FAQページとお問い合わせフォームへの外部リンクを提供する。
 * リンクは別タブで開く（target="_blank"）。
 *
 * @see {@link file://../../../docs/functions/settings/home.md} 詳細設計書
 */

// FAQ・問い合わせは外部リンクのみで状態管理が不要なため、Server Component として動作できるが、
// 親コンポーネント（SettingsPage）がClient Componentのため "use client" は不要。

const FAQ_URL = process.env.NEXT_PUBLIC_FAQ_URL ?? "https://sicou-lab.com/faq";
const CONTACT_URL = process.env.NEXT_PUBLIC_CONTACT_URL ?? "https://sicou-lab.com/contact";

/**
 * FAQ・問い合わせセクションコンポーネント
 *
 * FAQと問い合わせフォームへの外部リンクを提供する。
 *
 * @returns FAQ・問い合わせセクション
 */
export function FaqSection() {
  return (
    <section className="rounded-lg border border-gray-200 p-4 dark:border-gray-800 lg:p-6">
      <h2 className="text-xl font-bold text-[var(--color-text)]">FAQ・問い合わせ</h2>

      <div className="mt-4 flex flex-col gap-3 md:flex-row">
        <a
          href={FAQ_URL}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex min-h-[44px] items-center justify-center rounded-md border border-gray-300 px-4 py-2 text-lg text-[var(--color-text)] hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-800"
        >
          FAQを見る
        </a>
        <a
          href={CONTACT_URL}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex min-h-[44px] items-center justify-center rounded-md border border-gray-300 px-4 py-2 text-lg text-[var(--color-text)] hover:bg-gray-100 dark:border-gray-600 dark:hover:bg-gray-800"
        >
          お問い合わせ
        </a>
      </div>
    </section>
  );
}
