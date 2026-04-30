import Link from "next/link";
import { notFound } from "next/navigation";

import { RichTextRenderer } from "@/components/articles/RichTextRenderer";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ApiError } from "@/tools/api";
import { getArticle } from "@/services/articles";

export const revalidate = 60;

export async function generateStaticParams() {
  return [];
}

export const dynamicParams = true;

export async function generateMetadata({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  try {
    const article = await getArticle(id, { revalidate: 60 });
    return { title: article.title };
  } catch {
    return { title: "Article" };
  }
}

const dateFormatter = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "long",
  day: "2-digit",
});

export default async function ArticleDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  let article;
  try {
    article = await getArticle(id, { revalidate: 60 });
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) {
      notFound();
    }
    throw err;
  }

  const created = new Date(article.date_created);

  return (
    <article className="mx-auto w-full max-w-3xl px-4 py-12 sm:px-6">
      <Button asChild variant="link" className="mb-6 px-0">
        <Link href="/articles">← All articles</Link>
      </Button>

      <header className="space-y-3">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          {Number.isNaN(created.getTime()) ? "" : dateFormatter.format(created)}
        </p>
        <h1 className="text-4xl font-semibold leading-tight tracking-tight">
          {article.title}
        </h1>
      </header>

      <Separator className="my-8" />

      <RichTextRenderer html={article.article_text} />
    </article>
  );
}
