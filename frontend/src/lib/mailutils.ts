import type { Part } from "$lib/types";

export function fmtDate(iso: string): string {
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  return d.toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function fmtDateTime(iso: string): string {
  const d = new Date(iso);
  if (isNaN(d.getTime())) return iso;
  return d.toLocaleString(undefined, {
    weekday: "short",
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

// first text part in priority order, used for message list previews
export function textPreview(parts: Part[] | undefined): string {
  if (!parts) return "";
  for (const part of parts) {
    if (part.content_type === "text/plain" && part.data) return part.data;
  }
  for (const part of parts) {
    if (part.content_type === "text/html" && part.data) return part.data as string;
  }
  return "";
}
