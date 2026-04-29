import type { ClientLogEntry, LogLevel } from "@/types/api";

// Best-effort frontend logger. Sends entries to POST /api/logs on the
// browser; falls back to console.* on the server. Never throws — a logging
// failure must not cascade into another runtime error.

const ENDPOINT_PUBLIC =
  (process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost/api") + "/logs";

function safePost(entry: ClientLogEntry): void {
  if (typeof window === "undefined") {
    // On the server we just forward to stderr. The Go API container picks it
    // up indirectly through the actual SSR fetch error if any.
    // eslint-disable-next-line no-console
    console[entry.level === "info" ? "log" : entry.level](
      `[client-${entry.level}]`,
      entry.message,
      entry.meta ?? ""
    );
    return;
  }

  const payload = JSON.stringify({
    ...entry,
    user_agent: entry.user_agent ?? navigator.userAgent,
    url: entry.url ?? window.location.href,
  });

  try {
    if (typeof navigator !== "undefined" && navigator.sendBeacon) {
      const blob = new Blob([payload], { type: "application/json" });
      if (navigator.sendBeacon(ENDPOINT_PUBLIC, blob)) return;
    }
    void fetch(ENDPOINT_PUBLIC, {
      method: "POST",
      body: payload,
      headers: { "Content-Type": "application/json" },
      keepalive: true,
      credentials: "omit",
    }).catch(() => {});
  } catch {
    // Swallow.
  }
}

function buildEntry(
  level: LogLevel,
  message: unknown,
  meta?: Record<string, unknown>,
  stack?: string
): ClientLogEntry {
  let normalized: string;
  if (message instanceof Error) {
    normalized = message.message;
    stack = stack ?? message.stack;
  } else if (typeof message === "string") {
    normalized = message;
  } else {
    try {
      normalized = JSON.stringify(message);
    } catch {
      normalized = String(message);
    }
  }
  return { level, message: normalized, stack, meta };
}

export const logger = {
  info(message: unknown, meta?: Record<string, unknown>) {
    safePost(buildEntry("info", message, meta));
  },
  warn(message: unknown, meta?: Record<string, unknown>) {
    safePost(buildEntry("warn", message, meta));
  },
  error(message: unknown, meta?: Record<string, unknown>) {
    safePost(buildEntry("error", message, meta));
  },
};

let installed = false;

// Attaches window.onerror + unhandledrejection so uncaught browser errors are
// reported automatically. Idempotent.
export function installGlobalErrorReporter() {
  if (installed || typeof window === "undefined") return;
  installed = true;
  window.addEventListener("error", (event) => {
    safePost(
      buildEntry(
        "error",
        event.error ?? event.message,
        {
          filename: event.filename,
          lineno: event.lineno,
          colno: event.colno,
        },
        event.error?.stack
      )
    );
  });
  window.addEventListener("unhandledrejection", (event) => {
    const reason = event.reason;
    safePost(
      buildEntry(
        "error",
        reason instanceof Error ? reason : `Unhandled rejection: ${String(reason)}`,
        { kind: "unhandledrejection" },
        reason instanceof Error ? reason.stack : undefined
      )
    );
  });
}
