import Link from "next/link";
import { notFound } from "next/navigation";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { SmartImage } from "@/components/site/SmartImage";
import { ApiError } from "@/tools/api";
import { getProduct } from "@/services/products";

// ISR: revalidate detail pages on the same cadence as the listing.
export const revalidate = 60;

// Generate no static paths up-front — the products are seeded in bulk and
// pages render on first request, then get cached.
export async function generateStaticParams() {
  return [];
}

export const dynamicParams = true;

const TYPE_LABELS: Record<string, string> = {
  "10": "Drinks",
  "05": "Books",
  "20": "Electronics",
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  try {
    const product = await getProduct(id, { revalidate: 60 });
    return {
      title: product.product_name,
      description: `${product.product_name} — ${TYPE_LABELS[product.product_type] ?? "Product"}`,
    };
  } catch {
    return { title: "Product" };
  }
}

export default async function ProductDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let product;
  try {
    product = await getProduct(id, { revalidate: 60 });
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) {
      notFound();
    }
    throw err;
  }

  const price = Number(product.product_prices);
  const formatted = Number.isFinite(price)
    ? price.toLocaleString(undefined, { style: "currency", currency: "USD" })
    : `$${product.product_prices}`;
  const typeLabel =
    TYPE_LABELS[product.product_type] ?? `Type ${product.product_type}`;

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6">
      <Button asChild variant="link" className="mb-6 px-0">
        <Link href="/products">← Back to all products</Link>
      </Button>

      <div className="grid gap-10 md:grid-cols-2">
        <div className="relative aspect-square overflow-hidden rounded-2xl border bg-muted">
          <SmartImage
            imagePath={product.image_path}
            alt={product.product_name}
            fill
            sizes="(min-width: 768px) 50vw, 100vw"
            className="object-cover"
            priority
          />
        </div>
        <div className="flex flex-col gap-6">
          <div className="space-y-2">
            <Badge variant="secondary">{typeLabel}</Badge>
            <h1 className="text-3xl font-semibold tracking-tight">
              {product.product_name}
            </h1>
            <p className="text-2xl font-semibold">{formatted}</p>
          </div>

          <Separator />

          <dl className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <dt className="text-muted-foreground">SKU</dt>
              <dd className="font-mono">{product.product_id}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">In stock</dt>
              <dd>{product.product_quantity}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Type</dt>
              <dd>{typeLabel}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Added</dt>
              <dd>
                {new Date(product.created_at).toLocaleDateString(undefined, {
                  year: "numeric",
                  month: "short",
                  day: "2-digit",
                })}
              </dd>
            </div>
          </dl>

          <div className="flex flex-wrap gap-3">
            <Button size="lg" disabled>
              Add to cart (demo)
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link href="/contact">Ask a question</Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
