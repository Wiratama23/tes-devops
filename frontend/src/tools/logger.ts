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

// Handler refs are stored at module scope so uninstall can remove the exact
// same function references later. The "installed" flag lives on the window
// object so it survives module re-evaluation under Fast Refresh / Turbopack
// HMR — without that, a hot reload would reset a module-local flag while the
// previous listeners stayed attached, stacking handlers on every reload.
type ErrorReporterWindow = Window & {
  __tesdevopsErrorReporterInstalled?: boolean;
};

let errorHandler: ((event: ErrorEvent) => void) | null = null;
let rejectionHandler: ((event: PromiseRejectionEvent) => void) | null = null;

// Attaches window.onerror + unhandledrejection so uncaught browser errors are
// reported automatically. Idempotent across module reloads.
export function installGlobalErrorReporter() {
  if (typeof window === "undefined") return;
  const w = window as ErrorReporterWindow;
  if (w.__tesdevopsErrorReporterInstalled) return;
  w.__tesdevopsErrorReporterInstalled = true;

  errorHandler = (event: ErrorEvent) => {
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
  };

  rejectionHandler = (event: PromiseRejectionEvent) => {
    const reason = event.reason;
    safePost(
      buildEntry(
        "error",
        reason instanceof Error
          ? reason
          : `Unhandled rejection: ${String(reason)}`,
        { kind: "unhandledrejection" },
        reason instanceof Error ? reason.stack : undefined
      )
    );
  };

  window.addEventListener("error", errorHandler);
  window.addEventListener("unhandledrejection", rejectionHandler);
}

// Detaches the listeners installed by installGlobalErrorReporter. Safe to
// call when nothing was installed; safe to call repeatedly.
export function uninstallGlobalErrorReporter() {
  if (typeof window === "undefined") return;
  if (errorHandler) {
    window.removeEventListener("error", errorHandler);
    errorHandler = null;
  }
  if (rejectionHandler) {
    window.removeEventListener("unhandledrejection", rejectionHandler);
    rejectionHandler = null;
  }
  delete (window as ErrorReporterWindow).__tesdevopsErrorReporterInstalled;
}
