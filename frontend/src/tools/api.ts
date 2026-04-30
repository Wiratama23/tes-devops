// Typed API wrapper used by both Server Components (SSG/ISR) and the
// browser (CSR), built on top of axios. The base URL is selected based on
// whether the caller is on the server or client so the same code path works
// on both sides.
//
// Note: axios does not participate in Next.js' per-fetch ISR cache the way
// native fetch does. We control regeneration via page-level
// `export const revalidate = N` instead, which is enough for this app.

import axios, {
  AxiosError,
  type AxiosInstance,
  type AxiosRequestConfig,
} from "axios";

import { apiBaseUrlPublic, apiBaseUrlServer } from "@/tools/api-base";
import { ApiError, sanitizeErrorMessage } from "@/tools/api-error";

// Two singletons so server/client requests don't share interceptors. Each is
// created lazily on first use.
let serverClient: AxiosInstance | null = null;
let browserClient: AxiosInstance | null = null;

function getClient(): AxiosInstance {
  const onServer = typeof window === "undefined";
  if (onServer) {
    if (!serverClient) {
      serverClient = axios.create({
  baseURL: apiBaseUrlServer(),
        withCredentials: true,
        timeout: 30_000,
      });
    }
    return serverClient;
  }
  if (!browserClient) {
    browserClient = axios.create({
  baseURL: apiBaseUrlPublic(),
      withCredentials: true,
      timeout: 30_000,
    });
  }
  return browserClient;
}

export interface ApiOptions {
  method?: AxiosRequestConfig["method"];
  // JSON body. If you need raw FormData / a string, use `rawBody`.
  body?: unknown;
  rawBody?: BodyInit;
  headers?: Record<string, string>;
  signal?: AbortSignal;
  // Kept for backward compat with the previous fetch-based wrapper. They are
  // no-ops now — Next.js' per-fetch cache is a fetch-only feature. Use
  // page-level `export const revalidate = N` instead.
  revalidate?: number | false;
  tags?: string[];
}


function dataFromOptions(options: ApiOptions): unknown {
  if (options.rawBody !== undefined) return options.rawBody;
  if (options.body !== undefined) return options.body;
  return undefined;
}

function isJsonBody(options: ApiOptions): boolean {
  return options.rawBody === undefined && options.body !== undefined;
}

export async function apiFetch<T>(
  path: string,
  options: ApiOptions = {}
): Promise<T> {
  const { method = "GET", headers = {}, signal } = options;

  const config: AxiosRequestConfig = {
    url: path,
    method,
    headers: {
      Accept: "application/json",
      ...(isJsonBody(options) ? { "Content-Type": "application/json" } : {}),
      ...headers,
    },
    data: dataFromOptions(options),
    signal,
    // Receive raw bodies as text/blob and decide what to do per response.
    // Using `transformResponse: x => x` would lose axios' built-in JSON
    // parsing, so we let axios parse JSON normally and handle 204 below.
    validateStatus: () => true,
  };

  const client = getClient();

  let response;
  try {
    response = await client.request<T>(config);
  } catch (err) {
    // Axios only ends up here for network failures, timeouts, and aborts —
    // validateStatus above means non-2xx never throws.
    const ax = err as AxiosError;
    if (!ax) {
      throw new ApiError("Unknown error occurred", 0, "");
    }
    // ECONNREFUSED = Connection refused (backend not running)
    // ENOTFOUND = DNS resolution failed
    // ETIMEDOUT = Connection timeout
    const isNetworkError = ax.code === "ECONNREFUSED" || 
                           ax.code === "ENOTFOUND" || 
                           ax.code === "ETIMEDOUT";
    const message = ax.message || "network error";
    const sanitized = sanitizeErrorMessage(`${method} ${path} failed: ${message}`);
    throw new ApiError(sanitized, isNetworkError ? 0 : ax.response?.status || 0, "");
  }

  if (!response) {
    throw new ApiError("No response from server", 0, "");
  }

  if (response.status >= 200 && response.status < 300) {
    if (response.status === 204) {
      return undefined as T;
    }
    return response.data;
  }

  const bodyText =
    typeof response.data === "string"
      ? response.data
      : response.data
      ? (() => {
          try {
            return JSON.stringify(response.data);
          } catch {
            return String(response.data);
          }
        })()
      : "";

  const sanitized = sanitizeErrorMessage(`${method} ${path} failed: ${response.status}`);
  throw new ApiError(sanitized, response.status, bodyText);
}
