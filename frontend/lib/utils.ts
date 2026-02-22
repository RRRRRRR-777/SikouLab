import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Tailwindクラスをマージする。
 *
 * clsxで条件付きクラスを結合し、tailwind-mergeで競合するユーティリティクラスを上書き解決する。
 *
 * @param inputs - マージするクラス値
 * @returns マージされたクラス文字列
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
