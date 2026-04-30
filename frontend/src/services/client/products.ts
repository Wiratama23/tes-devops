import type { PaginatedProducts, Product } from "@/types/api";
import type { ProductCreateInput, ProductUpdateInput } from "@/schemas";
import { clientFetch } from "@/tools/client-api";

interface ListOptions {
  page?: number;
  signal?: AbortSignal;
}

export async function listProducts({
  page = 1,
  signal,
}: ListOptions = {}): Promise<PaginatedProducts> {
  return clientFetch<PaginatedProducts>(`/products?page=${page}`, { signal });
}

export async function getProduct(
  id: string,
  { signal }: { signal?: AbortSignal } = {}
): Promise<Product> {
  return clientFetch<Product>(`/products/${encodeURIComponent(id)}`, { signal });
}

export async function createProduct(
  input: ProductCreateInput & { created_by: string }
): Promise<Product> {
  return clientFetch<Product>("/products", {
    method: "POST",
    body: input,
  });
}

export async function updateProduct(
  id: string,
  input: ProductUpdateInput
): Promise<Product> {
  return clientFetch<Product>(`/products/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: input,
  });
}

export async function deleteProduct(id: string): Promise<void> {
  await clientFetch<void>(`/products/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}
