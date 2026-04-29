import { apiFetch, apiBaseUrl } from "@/tools/api";
import type { UploadResponse } from "@/types/api";

export async function uploadImage(file: File): Promise<UploadResponse> {
  const form = new FormData();
  form.append("file", file);
  return apiFetch<UploadResponse>("/uploads/images", {
    method: "POST",
    rawBody: form,
  });
}

// Resolves a backend image_path string (e.g. "assets/coffee.jpg") into a
// browser-fetchable URL (`/api/assets/coffee.jpg`). Returns the public
// fallback when the path is empty.
export function resolveImageUrl(path: string | null | undefined): string {
  if (!path) {
    return process.env.NEXT_PUBLIC_DEFAULT_IMAGE ?? "/assets/default_image.jpg";
  }
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }
  if (path.startsWith("/")) {
    return path;
  }
  // "assets/foo.jpg" -> "<api-base>/assets/foo.jpg"
  if (path.startsWith("assets/")) {
    return `${apiBaseUrl()}/${path}`;
  }
  return `${apiBaseUrl()}/assets/${path}`;
}
