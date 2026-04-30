import Link from "next/link";

import { ProductCard } from "@/components/products/ProductCard";
import { Button } from "@/components/ui/button";
import { listProducts } from "@/services/products";

// ISR: re-validate the products listing every 60s.
export const revalidate = 60;

interface SearchParams {
  page?: string;
}

export const metadata = {
  title: "Products",
  description: "Browse the full product catalog.",
};

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: Promise<SearchParams>;
}) {
  const { page: pageParam } = await searchParams;
  const page = Math.max(1, Number(pageParam) || 1);

  const data = await listProducts({ page, revalidate: 60 });
  const items = data.data ?? [];

  return (
    <div className="mx-auto w-full max-w-6xl space-y-8 px-4 py-12 sm:px-6">
      <header className="flex flex-col gap-2">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Products · ISR (revalidate=60)
        </p>
        <h1 className="text-4xl font-semibold tracking-tight">All products</h1>
        <p className="max-w-2xl text-muted-foreground">
          The full catalog, regenerated automatically when the admin makes
          changes.
        </p>
      </header>

      {items.length === 0 ? (
        <p className="rounded-lg border border-dashed p-10 text-center text-muted-foreground">
          No products found on page {page}.
        </p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {items.map((product, idx) => (
            <ProductCard
              key={product.product_id}
              product={product}
              priority={idx < 4}
            />
          ))}
        </div>
      )}

      <Pagination page={page} pageSize={data.limit ?? 10} hasMore={items.length === (data.limit ?? 10)} />
    </div>
  );
}

function Pagination({
  page,
  hasMore,
}: {
  page: number;
  pageSize: number;
  hasMore: boolean;
}) {
  return (
    <nav className="flex items-center justify-between border-t pt-6">
      <Button
        asChild={page > 1}
        variant="outline"
        size="sm"
        disabled={page <= 1}
      >
        {page > 1 ? (
          <Link href={{ pathname: "/products", query: { page: page - 1 } }}>
            ← Previous
          </Link>
        ) : (
          <span>← Previous</span>
        )}
      </Button>
      <span className="text-sm text-muted-foreground">Page {page}</span>
      <Button
        asChild={hasMore}
        variant="outline"
        size="sm"
        disabled={!hasMore}
      >
        {hasMore ? (
          <Link href={{ pathname: "/products", query: { page: page + 1 } }}>
            Next →
          </Link>
        ) : (
          <span>Next →</span>
        )}
      </Button>
    </nav>
  );
}
