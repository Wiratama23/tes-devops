"use client";

import { useQuery } from "@tanstack/react-query";
import { FileText, Package, Users } from "lucide-react";

import { ArticleCard } from "@/components/articles/ArticleCard";
import { ArticleCardSkeleton } from "@/components/articles/ArticleCardSkeleton";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { listArticles } from "@/services/articles";
import { listProducts } from "@/services/products";
import { listUsers } from "@/services/users";

export function DashboardStats() {
  const products = useQuery({
    queryKey: ["products", "page", 1],
    queryFn: ({ signal }) => listProducts({ page: 1, signal }),
    staleTime: 60_000,
  });
  const articles = useQuery({
    queryKey: ["articles", "page", 1],
    queryFn: ({ signal }) => listArticles({ page: 1, signal }),
    staleTime: 60_000,
  });
  const users = useQuery({
    queryKey: ["users"],
    queryFn: ({ signal }) => listUsers({ signal }),
    staleTime: 60_000,
  });

  return (
    <div className="space-y-8">
      <div className="grid gap-4 sm:grid-cols-3">
        <StatCard
          icon={<Package className="h-4 w-4" />}
          label="Products on first page"
          value={products.isLoading ? null : products.data?.data.length ?? 0}
        />
        <StatCard
          icon={<FileText className="h-4 w-4" />}
          label="Articles total"
          value={articles.isLoading ? null : articles.data?.total_count ?? 0}
        />
        <StatCard
          icon={<Users className="h-4 w-4" />}
          label="Users total"
          value={users.isLoading ? null : users.data?.length ?? 0}
        />
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold tracking-tight">
          Latest articles
        </h2>
        {articles.isLoading ? (
          <div className="grid gap-4 md:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <ArticleCardSkeleton key={i} />
            ))}
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-3">
            {(articles.data?.data ?? []).slice(0, 3).map((article) => (
              <ArticleCard key={article.articles_id} article={article} />
            ))}
            {articles.data && articles.data.data.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                No articles yet — head to{" "}
                <a className="underline" href="/admin/articles">
                  Articles
                </a>{" "}
                to create one.
              </p>
            ) : null}
          </div>
        )}
      </section>
    </div>
  );
}

function StatCard({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: number | null;
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{label}</CardTitle>
        <span className="text-muted-foreground">{icon}</span>
      </CardHeader>
      <CardContent>
        {value === null ? (
          <Skeleton className="h-7 w-16" />
        ) : (
          <p className="text-2xl font-semibold">{value.toLocaleString()}</p>
        )}
      </CardContent>
    </Card>
  );
}
