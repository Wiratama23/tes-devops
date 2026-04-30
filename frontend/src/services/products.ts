import { apiFetch } from "@/tools/api";
import type {
  PaginatedProducts,
  Product,
} from "@/types/api";
import type {
  ProductCreateInput,
  ProductUpdateInput,
} from "@/schemas";

interface ListOptions {
  page?: number;
  // Only the server-side calls (SSG/ISR) pass `revalidate`. CSR callers using
  // SWR rely on its cache instead.
  revalidate?: number | false;
  tag?: string;
  // Forwarded from SWR so cancelled queries actually abort the underlying
  // axios request instead of leaving it in-flight.
  signal?: AbortSignal;
}

export async function listProducts({
  page = 1,
  revalidate,
  tag = "products",
  signal,
}: ListOptions = {}): Promise<PaginatedProducts> {
  return apiFetch<PaginatedProducts>(
    `/products?page=${page}`,
    { revalidate, tags: [tag], signal }
  );
}

export async function getProduct(
  id: string,
  options: { revalidate?: number | false; signal?: AbortSignal } = {}
): Promise<Product> {
  return apiFetch<Product>(`/products/${encodeURIComponent(id)}`, {
    revalidate: options.revalidate,
    tags: [`product:${id}`],
    signal: options.signal,
  });
}

export async function createProduct(
  input: ProductCreateInput & { created_by: string }
): Promise<Product> {
  return apiFetch<Product>("/products", {
    method: "POST",
    body: input,
  });
}

export async function updateProduct(
  id: string,
  input: ProductUpdateInput
): Promise<Product> {
  return apiFetch<Product>(`/products/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: input,
  });
}

export async function deleteProduct(id: string): Promise<void> {
  await apiFetch<void>(`/products/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}
