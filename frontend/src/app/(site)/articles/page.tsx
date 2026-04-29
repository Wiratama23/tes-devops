import Link from "next/link";

import { ArticleCard } from "@/components/articles/ArticleCard";
import { Button } from "@/components/ui/button";
import { listArticles } from "@/services/articles";

export const revalidate = 60;

interface SearchParams {
  page?: string;
}

export const metadata = {
  title: "Articles",
  description: "Engineering articles published from the Tesdevops admin panel.",
};

export default async function ArticlesPage({
  searchParams,
}: {
  searchParams: Promise<SearchParams>;
}) {
  const { page: pageParam } = await searchParams;
  const page = Math.max(1, Number(pageParam) || 1);

  const data = await listArticles({ page, revalidate: 60 });
  const items = data.data ?? [];
  const limit = data.limit ?? 10;
  const totalPages = Math.max(1, Math.ceil((data.total_count ?? 0) / limit));

  return (
    <div className="mx-auto w-full max-w-5xl space-y-8 px-4 py-12 sm:px-6">
      <header className="space-y-2">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Articles · ISR (revalidate=60)
        </p>
        <h1 className="text-4xl font-semibold tracking-tight">All articles</h1>
        <p className="max-w-2xl text-muted-foreground">
          Long-form notes from the team. New posts are picked up at most a
          minute after the admin publishes them.
        </p>
      </header>

      {items.length === 0 ? (
        <p className="rounded-lg border border-dashed p-10 text-center text-muted-foreground">
          No articles on page {page}.
        </p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {items.map((article) => (
            <ArticleCard key={article.articles_id} article={article} />
          ))}
        </div>
      )}

      <Pagination page={page} totalPages={totalPages} />
    </div>
  );
}

function Pagination({ page, totalPages }: { page: number; totalPages: number }) {
  const hasPrev = page > 1;
  const hasNext = page < totalPages;
  return (
    <nav className="flex items-center justify-between border-t pt-6">
      <Button
        asChild={hasPrev}
        variant="outline"
        size="sm"
        disabled={!hasPrev}
      >
        {hasPrev ? (
          <Link href={{ pathname: "/articles", query: { page: page - 1 } }}>
            ← Previous
          </Link>
        ) : (
          <span>← Previous</span>
        )}
      </Button>
      <span className="text-sm text-muted-foreground">
        Page {page} / {totalPages}
      </span>
      <Button
        asChild={hasNext}
        variant="outline"
        size="sm"
        disabled={!hasNext}
      >
        {hasNext ? (
          <Link href={{ pathname: "/articles", query: { page: page + 1 } }}>
            Next →
          </Link>
        ) : (
          <span>Next →</span>
        )}
      </Button>
    </nav>
  );
}
