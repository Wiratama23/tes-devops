import Link from "next/link";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/tools/utils";
import type { Article } from "@/types/api";

interface ArticleCardProps {
  article: Article;
  className?: string;
}

const dateFormatter = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "2-digit",
});

function stripHtml(value: string): string {
  return value
    .replace(/<[^>]+>/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}

export function ArticleCard({ article, className }: ArticleCardProps) {
  const created = new Date(article.date_created);
  const preview = stripHtml(article.article_text).slice(0, 180);

  return (
    <Card
      className={cn(
        "h-full transition-all hover:-translate-y-0.5 hover:shadow-md",
        className
      )}
    >
      <Link
        href={`/articles/${article.articles_id}`}
        className="block h-full focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        prefetch
      >
        <CardHeader>
          <p className="text-xs uppercase tracking-wide text-muted-foreground">
            {Number.isNaN(created.getTime())
              ? ""
              : dateFormatter.format(created)}
          </p>
          <CardTitle className="line-clamp-2 text-lg">{article.title}</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="line-clamp-3 text-sm text-muted-foreground">
            {preview}
            {preview.length === 180 ? "…" : ""}
          </p>
        </CardContent>
      </Link>
    </Card>
  );
}
