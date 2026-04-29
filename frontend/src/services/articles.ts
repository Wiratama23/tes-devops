import { apiFetch } from "@/tools/api";
import type { Article, PaginatedArticles } from "@/types/api";
import type {
  ArticleCreateInput,
  ArticleUpdateInput,
} from "@/schemas";

interface ListOptions {
  page?: number;
  revalidate?: number | false;
  tag?: string;
}

export async function listArticles({
  page = 1,
  revalidate,
  tag = "articles",
}: ListOptions = {}): Promise<PaginatedArticles> {
  return apiFetch<PaginatedArticles>(`/articles?page=${page}`, {
    revalidate,
    tags: [tag],
  });
}

export async function getArticle(
  id: number | string,
  options: { revalidate?: number | false } = {}
): Promise<Article> {
  return apiFetch<Article>(`/articles/${id}`, {
    revalidate: options.revalidate,
    tags: [`article:${id}`],
  });
}

export async function createArticle(
  input: ArticleCreateInput & { uid: string }
): Promise<Article> {
  return apiFetch<Article>("/articles", {
    method: "POST",
    body: input,
  });
}

export async function updateArticle(
  id: number | string,
  input: ArticleUpdateInput
): Promise<Article> {
  return apiFetch<Article>(`/articles/${id}`, {
    method: "PUT",
    body: input,
  });
}

export async function deleteArticle(id: number | string): Promise<void> {
  await apiFetch<void>(`/articles/${id}`, { method: "DELETE" });
}
