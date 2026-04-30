import type { AuthUser, LoginResponse } from "@/types/api";
import { clientFetch } from "@/tools/client-api";

export async function login(input: {
  username: string;
  password: string;
}): Promise<LoginResponse> {
  return clientFetch<LoginResponse>("/auth/login", {
    method: "POST",
    body: input,
  });
}

export async function logout(): Promise<void> {
  await clientFetch<void>("/auth/logout", { method: "POST" });
}

export async function me({
  signal,
}: { signal?: AbortSignal } = {}): Promise<AuthUser> {
  return clientFetch<AuthUser>("/auth/me", { method: "GET", signal });
}

export async function refreshToken(): Promise<LoginResponse> {
  return clientFetch<LoginResponse>("/auth/refresh", { method: "POST" });
}
