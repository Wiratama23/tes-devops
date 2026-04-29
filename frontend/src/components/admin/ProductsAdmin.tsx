"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil, Plus, Trash2 } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";
import { toast } from "sonner";

import { ConfirmDialog } from "@/components/admin/ConfirmDialog";
import { ProductDialog } from "@/components/admin/ProductDialog";
import { SmartImage } from "@/components/site/SmartImage";
import { Badge } from "@/components/ui/badge";
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
  createProduct,
  deleteProduct,
  listProducts,
  updateProduct,
} from "@/services/products";
import type { PaginatedProducts, Product } from "@/types/api";

export function ProductsAdmin() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState<Product | undefined>(undefined);
  const [creating, setCreating] = useState(false);
  const [confirming, setConfirming] = useState<Product | null>(null);

  const meQuery = useQuery({
    queryKey: ["auth", "me"],
    queryFn: me,
    staleTime: Infinity,
  });

  const productsKey = useMemo(() => ["products", "page", page] as const, [page]);

  const productsQuery = useQuery({
    queryKey: productsKey,
    queryFn: () => listProducts({ page }),
  });

  const createMutation = useMutation({
    mutationFn: createProduct,
    onSuccess: () => {
      toast.success("Product created");
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
    onError: (err) => {
      logger.error("create product failed", { kind: "admin.product.create", err: String(err) });
      toast.error("Failed to create product.");
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      values,
    }: {
      id: string;
      values: Parameters<typeof updateProduct>[1];
    }) => updateProduct(id, values),
    onSuccess: () => {
      toast.success("Product updated");
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
    onError: (err) => {
      logger.error("update product failed", { kind: "admin.product.update", err: String(err) });
      toast.error("Failed to update product.");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProduct,
    onMutate: async (id: string) => {
      await queryClient.cancelQueries({ queryKey: ["products"] });
      const snapshot = queryClient.getQueryData<PaginatedProducts>(productsKey);
      if (snapshot) {
        queryClient.setQueryData<PaginatedProducts>(productsKey, {
          ...snapshot,
          data: snapshot.data.filter((p) => p.product_id !== id),
        });
      }
      return { snapshot };
    },
    onError: (err, _id, ctx) => {
      logger.error("delete product failed", { kind: "admin.product.delete", err: String(err) });
      if (ctx?.snapshot) {
        queryClient.setQueryData(productsKey, ctx.snapshot);
      }
      toast.error("Couldn't delete that product. Restored it.");
    },
    onSuccess: () => {
      toast.success("Product deleted");
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  const data = productsQuery.data?.data ?? [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-end">
        <Button onClick={() => setCreating(true)} disabled={!meQuery.data}>
          <Plus className="mr-2 h-4 w-4" />
          New product
        </Button>
      </div>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-16"></TableHead>
              <TableHead>Name</TableHead>
              <TableHead>SKU</TableHead>
              <TableHead className="hidden md:table-cell">Type</TableHead>
              <TableHead className="hidden md:table-cell">Stock</TableHead>
              <TableHead>Price</TableHead>
              <TableHead className="w-32 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {productsQuery.isLoading
              ? Array.from({ length: 6 }).map((_, i) => (
                  <TableRow key={`s-${i}`}>
                    <TableCell>
                      <Skeleton className="h-12 w-12 rounded" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-40" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-20" />
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      <Skeleton className="h-4 w-16" />
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      <Skeleton className="h-4 w-12" />
                    </TableCell>
                    <TableCell>
                      <Skeleton className="h-4 w-16" />
                    </TableCell>
                    <TableCell />
                  </TableRow>
                ))
              : data.map((product) => (
                  <TableRow
                    key={product.product_id}
                    data-testid={`product-row-${product.product_id}`}
                  >
                    <TableCell>
                      <div className="relative h-12 w-12 overflow-hidden rounded border bg-muted">
                        <SmartImage
                          imagePath={product.image_path}
                          alt={product.product_name}
                          fill
                          sizes="48px"
                          className="object-cover"
                        />
                      </div>
                    </TableCell>
                    <TableCell>
                      <Link
                        href={`/products/${encodeURIComponent(product.product_id)}`}
                        className="font-medium hover:underline"
                      >
                        {product.product_name}
                      </Link>
                    </TableCell>
                    <TableCell className="font-mono text-xs">
                      {product.product_id}
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      <Badge variant="secondary">{product.product_type}</Badge>
                    </TableCell>
                    <TableCell className="hidden md:table-cell">
                      {product.product_quantity}
                    </TableCell>
                    <TableCell>${product.product_prices}</TableCell>
                    <TableCell className="text-right">
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => setEditing(product)}
                        aria-label={`Edit ${product.product_name}`}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => setConfirming(product)}
                        aria-label={`Delete ${product.product_name}`}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}

            {!productsQuery.isLoading && data.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="py-10 text-center text-muted-foreground">
                  No products on page {page}.
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
          disabled={page <= 1 || productsQuery.isFetching}
          onClick={() => setPage((p) => Math.max(1, p - 1))}
        >
          ← Previous
        </Button>
        <span className="text-sm text-muted-foreground">Page {page}</span>
        <Button
          variant="outline"
          size="sm"
          disabled={data.length < (productsQuery.data?.limit ?? 10) || productsQuery.isFetching}
          onClick={() => setPage((p) => p + 1)}
        >
          Next →
        </Button>
      </div>

      <ProductDialog
        open={creating || Boolean(editing)}
        onOpenChange={(open) => {
          if (!open) {
            setCreating(false);
            setEditing(undefined);
          }
        }}
        product={editing}
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
        title="Delete product?"
        description={`This will remove "${confirming?.product_name ?? ""}" from the catalog.`}
        destructive
        confirmLabel="Delete"
        busy={deleteMutation.isPending}
        onConfirm={() => {
          if (!confirming) return;
          deleteMutation.mutate(confirming.product_id);
          setConfirming(null);
        }}
      />
    </div>
  );
}
