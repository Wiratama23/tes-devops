import { apiFetch } from "@/tools/api";
import type { User } from "@/types/api";

export async function listUsers(): Promise<User[]> {
  return apiFetch<User[]>("/users", { revalidate: 60 });
}

export async function getUser(uid: string): Promise<User> {
  return apiFetch<User>(`/users/${encodeURIComponent(uid)}`);
}
