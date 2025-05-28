import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | undefined | null): string {
  if (!date) return 'N/A';
  try {
    return new Date(date).toLocaleString();
  } catch {
    return 'Invalid date';
  }
}
