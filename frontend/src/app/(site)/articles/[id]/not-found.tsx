import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function ArticleNotFound() {
  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col items-center gap-4 px-4 py-24 text-center sm:px-6">
      <p className="text-sm uppercase tracking-wide text-muted-foreground">
        404
      </p>
      <h1 className="text-3xl font-semibold tracking-tight">
        Article not found
      </h1>
      <p className="text-muted-foreground">
        We can&apos;t find that article. It may have been removed.
      </p>
      <Button asChild>
        <Link href="/articles">Back to articles</Link>
      </Button>
    </div>
  );
}
