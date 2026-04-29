import { apiFetch } from "@/tools/api";
import type { User } from "@/types/api";

interface ListOptions {
  signal?: AbortSignal;
}

export async function listUsers({ signal }: ListOptions = {}): Promise<User[]> {
  return apiFetch<User[]>("/users", { revalidate: 60, signal });
}

export async function getUser(
  uid: string,
  { signal }: ListOptions = {}
): Promise<User> {
  return apiFetch<User>(`/users/${encodeURIComponent(uid)}`, { signal });
}
