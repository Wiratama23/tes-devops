import type { User } from "@/types/api";
import { clientFetch } from "@/tools/client-api";

interface ListOptions {
  signal?: AbortSignal;
}

export async function listUsers({ signal }: ListOptions = {}): Promise<User[]> {
  return clientFetch<User[]>("/users", { signal });
}

export async function getUser(
  uid: string,
  { signal }: ListOptions = {}
): Promise<User> {
  return clientFetch<User>(`/users/${encodeURIComponent(uid)}`, { signal });
}
