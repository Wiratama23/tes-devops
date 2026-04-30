import Link from "next/link";
import { Suspense } from "react";

import { ArticleCard } from "@/components/articles/ArticleCard";
import { ArticleCardSkeleton } from "@/components/articles/ArticleCardSkeleton";
import { ProductCard } from "@/components/products/ProductCard";
import { ProductCardSkeleton } from "@/components/products/ProductCardSkeleton";
import { Button } from "@/components/ui/button";
import { listArticles } from "@/services/articles";
import { listProducts } from "@/services/products";

// Revalidate every 10 minutes. Falls back to stale content if backend unavailable during build.
export const revalidate = 600;

export const metadata = {
  title: "Home",
};

async function FeaturedProducts() {
  try {
    const data = await listProducts({ revalidate: 600 });
    const items = data.data.slice(0, 4);
    if (items.length === 0) {
      return (
        <p className="text-sm text-muted-foreground">No products available yet.</p>
      );
    }
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {items.map((product, idx) => (
          <ProductCard
            key={product.product_id}
            product={product}
            priority={idx === 0}
          />
        ))}
      </div>
    );
  } catch (error) {
    // If backend is unavailable during build, render empty state
    return (
      <p className="text-sm text-muted-foreground">No products available yet.</p>
    );
  }
}

async function LatestArticles() {
  try {
    const data = await listArticles({ revalidate: 600 });
    const items = data.data.slice(0, 3);
    if (items.length === 0) {
      return (
        <p className="text-sm text-muted-foreground">
          No articles published yet.
        </p>
      );
    }
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {items.map((article) => (
          <ArticleCard key={article.articles_id} article={article} />
        ))}
      </div>
    );
  } catch (error) {
    // If backend is unavailable during build, render empty state
    return (
      <p className="text-sm text-muted-foreground">
        No articles published yet.
      </p>
    );
  }
}

export default function HomePage() {
  return (
    <div className="mx-auto w-full max-w-6xl space-y-20 px-4 py-12 sm:px-6 sm:py-16">
      <section className="grid gap-10 md:grid-cols-2 md:items-center md:gap-16">
        <div className="space-y-6 [animation:var(--animate-slide-up)]">
          <span className="inline-flex items-center rounded-full border border-border/60 bg-muted/40 px-3 py-1 text-xs uppercase tracking-wide text-muted-foreground">
            Next.js 16 · Go 1.26 · Postgres 16
          </span>
          <h1 className="text-4xl font-semibold leading-tight tracking-tight sm:text-5xl">
            A clean, fast storefront with a Go-powered admin behind the curtain.
          </h1>
          <p className="text-lg leading-relaxed text-muted-foreground">
            Browse curated products, read deep-dive articles, and ship updates
            instantly with incremental static regeneration. The whole stack is
            wired together with bun, Tailwind v4, and a typed Go API.
          </p>
          <div className="flex flex-wrap gap-3">
            <Button asChild size="lg">
              <Link href="/products">Browse products</Link>
            </Button>
            <Button asChild size="lg" variant="outline">
              <Link href="/articles">Read articles</Link>
            </Button>
          </div>
        </div>
        <div className="relative aspect-[4/3] overflow-hidden rounded-2xl border bg-gradient-to-br from-primary/15 via-background to-muted/40 shadow-sm">
          <div className="absolute inset-0 grid grid-cols-3 grid-rows-3 gap-2 p-6">
            {Array.from({ length: 9 }).map((_, i) => (
              <div
                key={i}
                className="rounded-md border border-border/40 bg-background/40 backdrop-blur-sm transition-transform hover:scale-[1.02]"
                style={{
                  animation: `var(--animate-fade-in)`,
                  animationDelay: `${i * 60}ms`,
                  animationFillMode: "backwards",
                }}
              />
            ))}
          </div>
        </div>
      </section>

      <section className="space-y-6">
        <div className="flex items-end justify-between">
          <div>
            <h2 className="text-2xl font-semibold tracking-tight">
              Featured products
            </h2>
            <p className="text-sm text-muted-foreground">
              A handful of items pulled live from the Go API at build time.
            </p>
          </div>
          <Button asChild variant="link" className="hidden sm:inline-flex">
            <Link href="/products">View all →</Link>
          </Button>
        </div>
        <Suspense fallback={<ProductsSkeletonRow />}>
          <FeaturedProducts />
        </Suspense>
      </section>

      <section className="space-y-6">
        <div className="flex items-end justify-between">
          <div>
            <h2 className="text-2xl font-semibold tracking-tight">
              Latest articles
            </h2>
            <p className="text-sm text-muted-foreground">
              Engineering notes and updates from the team.
            </p>
          </div>
          <Button asChild variant="link" className="hidden sm:inline-flex">
            <Link href="/articles">View all →</Link>
          </Button>
        </div>
        <Suspense fallback={<ArticlesSkeletonRow />}>
          <LatestArticles />
        </Suspense>
      </section>
    </div>
  );
}

function ProductsSkeletonRow() {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {Array.from({ length: 4 }).map((_, i) => (
        <ProductCardSkeleton key={i} />
      ))}
    </div>
  );
}

function ArticlesSkeletonRow() {
  return (
    <div className="grid gap-4 md:grid-cols-3">
      {Array.from({ length: 3 }).map((_, i) => (
        <ArticleCardSkeleton key={i} />
      ))}
    </div>
  );
}