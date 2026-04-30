import { apiBaseUrlPublic } from "@/tools/api-base";
import { ApiError, sanitizeErrorMessage } from "@/tools/api-error";

export interface ClientApiOptions {
  method?: string;
  body?: unknown;
  rawBody?: BodyInit;
  headers?: Record<string, string>;
  signal?: AbortSignal;
}

function dataFromOptions(options: ClientApiOptions): BodyInit | undefined {
  if (options.rawBody !== undefined) return options.rawBody;
  if (options.body !== undefined) return JSON.stringify(options.body);
  return undefined;
}

function isJsonBody(options: ClientApiOptions): boolean {
  return options.rawBody === undefined && options.body !== undefined;
}

export async function clientFetch<T>(
  path: string,
  options: ClientApiOptions = {}
): Promise<T> {
  const { method = "GET", headers = {}, signal } = options;

  const url = `${apiBaseUrlPublic()}${path}`;

  let response: Response;
  try {
    response = await fetch(url, {
      method,
      headers: {
        Accept: "application/json",
        ...(isJsonBody(options) ? { "Content-Type": "application/json" } : {}),
        ...headers,
      },
      body: dataFromOptions(options),
      signal,
      credentials: "include",
    });
  } catch (err) {
    const message = err instanceof Error ? err.message : "network error";
    const sanitized = sanitizeErrorMessage(`${method} ${path} failed: ${message}`);
    throw new ApiError(sanitized, 0, "");
  }

  if (response.status >= 200 && response.status < 300) {
    if (response.status === 204) {
      return undefined as T;
    }
    const contentType = response.headers.get("content-type");
    if (contentType?.includes("application/json")) {
      return (await response.json()) as T;
    }
    return (await response.text()) as T;
  }

  const bodyText = await response.text();
  const sanitized = sanitizeErrorMessage(
    `${method} ${path} failed: ${response.status}`
  );
  throw new ApiError(sanitized, response.status, bodyText);
}
