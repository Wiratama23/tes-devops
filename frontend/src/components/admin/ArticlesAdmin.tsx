"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil, Plus, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";

import { ArticleDialog } from "@/components/admin/ArticleDialog";
import { ConfirmDialog } from "@/components/admin/ConfirmDialog";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { logger } from "@/tools/logger";
import { me } from "@/services/auth";
import {
  createArticle,
  deleteArticle,
  listArticles,
  updateArticle,
} from "@/services/articles";
import type { Article, PaginatedArticles } from "@/types/api";

function stripHtml(value: string): string {
  return value.replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
}

export function ArticlesAdmin() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState<Article | undefined>(undefined);
  const [creating, setCreating] = useState(false);
  const [confirming, setConfirming] = useState<Article | null>(null);

  const meQuery = useQuery({
    queryKey: ["auth", "me"],
    queryFn: me,
    staleTime: Infinity,
  });

  const articlesKey = useMemo(() => ["articles", "page", page] as const, [page]);

  const articlesQuery = useQuery({
    queryKey: articlesKey,
    queryFn: () => listArticles({ page }),
  });

  const createMutation = useMutation({
    mutationFn: createArticle,
    onSuccess: () => {
      toast.success("Article created");
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
    onError: (err) => {
      logger.error("create article failed", { kind: "admin.article.create", err: String(err) });
      toast.error("Failed to create article.");
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      values,
    }: {
      id: number;
      values: Parameters<typeof updateArticle>[1];
    }) => updateArticle(id, values),
    onSuccess: () => {
      toast.success("Article updated");
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
    onError: (err) => {
      logger.error("update article failed", { kind: "admin.article.update", err: String(err) });
      toast.error("Failed to update article.");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteArticle,
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ["articles"] });
      const snapshot =
        queryClient.getQueryData<PaginatedArticles>(articlesKey);
      if (snapshot) {
        queryClient.setQueryData<PaginatedArticles>(articlesKey, {
          ...snapshot,
          data: snapshot.data.filter((a) => a.articles_id !== Number(id)),
        });
      }
      return { snapshot };
    },
    onError: (err, _id, ctx) => {
      logger.error("delete article failed", {
        kind: "admin.article.delete",
        err: String(err),
      });
      if (ctx?.snapshot) {
        queryClient.setQueryData(articlesKey, ctx.snapshot);
      }
      toast.error("Couldn't delete that article. Restored it.");
    },
    onSuccess: () => {
      toast.success("Article deleted");
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });

  const data = articlesQuery.data?.data ?? [];
  const limit = articlesQuery.data?.limit ?? 10;
  const totalPages = Math.max(
    1,
    Math.ceil((articlesQuery.data?.total_count ?? 0) / limit)
  );

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-end">
        <Button onClick={() => setCreating(true)} disabled={!meQuery.data}>
          <Plus className="mr-2 h-4 w-4" />
          New article
        </Button>
      </div>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Title</TableHead>
              <TableHead className="hidden md:table-cell">Preview</TableHead>
              <TableHead className="hidden lg:table-cell">Created</TableHead>
              <TableHead className="w-32 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {articlesQuery.isLoading
              ? Array.from({ length: 6 }).map((_, i) => (
                  <TableRow key={`s-${i}`}>
                    <TableCell>
                      <Skeleton className="h-4 w-48" />
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      <Skeleton className="h-4 w-full" />
                    </TableCell>
                    <TableCell className="hidden lg:table-cell">
                      <Skeleton className="h-4 w-24" />
                    </TableCell>
                    <TableCell />
                  </TableRow>
                ))
              : data.map((article) => (
                  <TableRow
                    key={article.articles_id}
                    data-testid={`article-row-${article.articles_id}`}
                  >
                    <TableCell className="font-medium">
                      {article.title}
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      <p className="line-clamp-2 max-w-md text-sm text-muted-foreground">
                        {stripHtml(article.article_text).slice(0, 200)}
                      </p>
                    </TableCell>
                    <TableCell className="hidden lg:table-cell text-xs text-muted-foreground">
                      {new Date(article.date_created).toLocaleDateString()}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => setEditing(article)}
                        aria-label={`Edit ${article.title}`}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => setConfirming(article)}
                        aria-label={`Delete ${article.title}`}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}

            {!articlesQuery.isLoading && data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="py-10 text-center text-muted-foreground">
                  No articles on page {page}.
                </TableCell>
              </TableRow>
            ) : null}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between">
        <Button
          variant="outline"
          size="sm"
          disabled={page <= 1 || articlesQuery.isFetching}
          onClick={() => setPage((p) => Math.max(1, p - 1))}
        >
          ← Previous
        </Button>
        <span className="text-sm text-muted-foreground">
          Page {page} / {totalPages}
        </span>
        <Button
          variant="outline"
          size="sm"
          disabled={page >= totalPages || articlesQuery.isFetching}
          onClick={() => setPage((p) => p + 1)}
        >
          Next →
        </Button>
      </div>

      <ArticleDialog
        open={creating || Boolean(editing)}
        onOpenChange={(open) => {
          if (!open) {
            setCreating(false);
            setEditing(undefined);
          }
        }}
        article={editing}
        currentUserId={meQuery.data?.uid ?? ""}
        onSubmit={async (payload) => {
          if (payload.mode === "create") {
            await createMutation.mutateAsync(payload.values);
          } else {
            await updateMutation.mutateAsync({
              id: payload.id,
              values: payload.values,
            });
          }
        }}
      />

      <ConfirmDialog
        open={Boolean(confirming)}
        onOpenChange={(open) => !open && setConfirming(null)}
        title="Delete article?"
        description={`This will remove "${confirming?.title ?? ""}" from the site.`}
        destructive
        confirmLabel="Delete"
        busy={deleteMutation.isPending}
        onConfirm={() => {
          if (!confirming) return;
          deleteMutation.mutate(confirming.articles_id);
          setConfirming(null);
        }}
      />
    </div>
  );
}
