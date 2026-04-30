import { apiFetch } from "@/tools/api";
import type { AuthUser, LoginResponse } from "@/types/api";

export async function login(input: {
  username: string;
  password: string;
}): Promise<LoginResponse> {
  return apiFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: input,
  });
}

export async function logout(): Promise<void> {
  await apiFetch<void>("/auth/logout", { method: "POST" });
}

export async function me({
  signal,
}: { signal?: AbortSignal } = {}): Promise<AuthUser> {
  return apiFetch<AuthUser>("/auth/me", { method: "GET", signal });
}

export async function refreshToken(): Promise<LoginResponse> {
  return apiFetch<LoginResponse>("/auth/refresh", { method: "POST" });
}
