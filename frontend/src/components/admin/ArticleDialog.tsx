"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";

import { TiptapEditor } from "@/components/articles/TiptapEditor";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  articleCreateSchema,
  articleUpdateSchema,
  type ArticleCreateInput,
  type ArticleUpdateInput,
} from "@/schemas";
import type { Article } from "@/types/api";

export interface ArticleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  article?: Article;
  currentUserId: string;
  onSubmit: (
    payload:
      | { mode: "create"; values: ArticleCreateInput & { uid: string } }
      | { mode: "update"; id: number; values: ArticleUpdateInput }
  ) => Promise<void>;
}

const defaultValues: ArticleCreateInput = {
  title: "",
  article_text: "",
};

export function ArticleDialog({
  open,
  onOpenChange,
  article,
  currentUserId,
  onSubmit,
}: ArticleDialogProps) {
  const isEdit = Boolean(article);
  const schema = isEdit ? articleUpdateSchema : articleCreateSchema;

  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ArticleCreateInput>({
    resolver: zodResolver(schema as never),
    defaultValues,
  });

  useEffect(() => {
    if (open) {
      reset(
        article
          ? { title: article.title, article_text: article.article_text }
          : defaultValues
      );
    }
  }, [open, article, reset]);

  async function handleFormSubmit(values: ArticleCreateInput) {
    if (isEdit && article) {
      await onSubmit({
        mode: "update",
        id: article.articles_id,
        values: {
          title: values.title,
          article_text: values.article_text,
        },
      });
    } else {
      await onSubmit({
        mode: "create",
        values: {
          ...values,
          uid: currentUserId,
        },
      });
    }
    onOpenChange(false);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit article" : "New article"}</DialogTitle>
          <DialogDescription>
            Title is plain text; the body uses a Tiptap rich-text editor and is
            stored as HTML on the backend.
          </DialogDescription>
        </DialogHeader>

        <form
          id="article-form"
          onSubmit={handleSubmit(handleFormSubmit)}
          className="space-y-4"
        >
          <div className="space-y-2">
            <Label htmlFor="article-title">Title</Label>
            <Input id="article-title" {...register("title")} />
            {errors.title ? (
              <p className="text-xs text-destructive">{errors.title.message}</p>
            ) : null}
          </div>

          <div className="space-y-2">
            <Label>Body</Label>
            <Controller
              control={control}
              name="article_text"
              render={({ field }) => (
                <TiptapEditor
                  value={field.value}
                  onChange={field.onChange}
                  disabled={isSubmitting}
                />
              )}
            />
            {errors.article_text ? (
              <p className="text-xs text-destructive">
                {errors.article_text.message}
              </p>
            ) : null}
          </div>
        </form>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSubmitting}
          >
            Cancel
          </Button>
          <Button form="article-form" type="submit" disabled={isSubmitting}>
            {isSubmitting
              ? "Saving…"
              : isEdit
              ? "Save changes"
              : "Publish article"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
