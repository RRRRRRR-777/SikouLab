import { cn } from "@/lib/utils";

/**
 * ローディング時のプレースホルダーコンポーネント。
 *
 * @param root0 - コンポーネントのプロパティ
 * @param root0.className - 追加のCSSクラス
 * @returns スケルトン要素
 */
function Skeleton({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="skeleton"
      className={cn("bg-accent animate-pulse rounded-md", className)}
      {...props}
    />
  );
}

export { Skeleton };
