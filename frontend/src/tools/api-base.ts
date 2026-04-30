const PUBLIC_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost/api";
const SERVER_BASE =
  process.env.INTERNAL_API_URL ??
  process.env.NEXT_PUBLIC_API_BASE_URL ??
  "http://localhost/api";

export function apiBaseUrl(): string {
  return typeof window === "undefined" ? SERVER_BASE : PUBLIC_BASE;
}

export function apiBaseUrlServer(): string {
  return SERVER_BASE;
}

export function apiBaseUrlPublic(): string {
  return PUBLIC_BASE;
}
