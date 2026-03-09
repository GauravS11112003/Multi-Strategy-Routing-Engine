import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs) {
  return twMerge(clsx(inputs));
}

export function formatNumber(num, decimals = 1) {
  if (num === null || num === undefined) return "—";
  if (num >= 1000) return (num / 1000).toFixed(1) + "k";
  return Number(num).toFixed(decimals);
}

export function formatDuration(minutes) {
  if (!minutes) return "—";
  const h = Math.floor(minutes / 60);
  const m = Math.round(minutes % 60);
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}

export const ROUTE_COLORS = [
  "#00C389",
  "#3b82f6",
  "#a855f7",
  "#f43f5e",
  "#f59e0b",
  "#06b6d4",
  "#ec4899",
  "#84cc16",
];
