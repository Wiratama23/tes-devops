import type { Article, PaginatedArticles } from "@/types/api";
import type { ArticleCreateInput, ArticleUpdateInput } from "@/schemas";
import { clientFetch } from "@/tools/client-api";

interface ListOptions {
  page?: number;
  signal?: AbortSignal;
}

export async function listArticles({
  page = 1,
  signal,
}: ListOptions = {}): Promise<PaginatedArticles> {
  return clientFetch<PaginatedArticles>(`/articles?page=${page}`, { signal });
}

export async function getArticle(
  id: number | string,
  { signal }: { signal?: AbortSignal } = {}
): Promise<Article> {
  return clientFetch<Article>(`/articles/${id}`, { signal });
}

export async function createArticle(
  input: ArticleCreateInput & { uid: string }
): Promise<Article> {
  return clientFetch<Article>("/articles", {
    method: "POST",
    body: input,
  });
}

export async function updateArticle(
  id: number | string,
  input: ArticleUpdateInput
): Promise<Article> {
  return clientFetch<Article>(`/articles/${id}`, {
    method: "PUT",
    body: input,
  });
}

export async function deleteArticle(id: number | string): Promise<void> {
  await clientFetch<void>(`/articles/${id}`, { method: "DELETE" });
}
