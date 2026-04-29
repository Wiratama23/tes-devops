"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";

import { ImageDropzone } from "@/components/products/ImageDropzone";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { logger } from "@/tools/logger";
import {
  productCreateSchema,
  productUpdateSchema,
  type ProductCreateInput,
  type ProductUpdateInput,
} from "@/schemas";
import { uploadImage } from "@/services/uploads";
import { resolveImageUrl } from "@/services/uploads";
import type { Product } from "@/types/api";

const TYPE_OPTIONS = [
  { value: "10", label: "Drinks" },
  { value: "05", label: "Books" },
  { value: "20", label: "Electronics" },
];

export interface ProductDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  // When `product` is undefined we render the create form, otherwise the
  // edit form. Submission funnels through the same callback so the page can
  // wire it up to a mutation.
  product?: Product;
  onSubmit: (
    values:
      | { mode: "create"; values: ProductCreateInput & { created_by: string } }
      | { mode: "update"; id: string; values: ProductUpdateInput }
  ) => Promise<void>;
  currentUserId: string;
}

type FormValues = ProductCreateInput;

const defaultValues: FormValues = {
  product_id: "",
  product_name: "",
  product_quantity: 0,
  product_prices: "",
  product_type: "10",
  image_path: "",
};

export function ProductDialog({
  open,
  onOpenChange,
  product,
  onSubmit,
  currentUserId,
}: ProductDialogProps) {
  const isEdit = Boolean(product);
  const schema = isEdit ? productUpdateSchema : productCreateSchema;

  const {
    register,
    control,
    handleSubmit,
    setValue,
    reset,
    formState: { errors, isSubmitting },
    watch,
  } = useForm<FormValues>({
    resolver: zodResolver(schema as never),
    defaultValues,
  });

  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    if (open) {
      if (product) {
        reset({
          product_id: product.product_id,
          product_name: product.product_name,
          product_quantity: product.product_quantity,
          product_prices: String(product.product_prices),
          product_type: product.product_type,
          image_path: product.image_path,
        });
      } else {
        reset(defaultValues);
      }
      setFile(null);
    }
  }, [open, product, reset]);

  // eslint-disable-next-line react-hooks/incompatible-library
  const imagePathValue = watch("image_path");

  async function handleFormSubmit(values: FormValues) {
    let finalImagePath = values.image_path;
    if (file) {
      try {
        setUploading(true);
        const result = await uploadImage(file);
        finalImagePath = result.path;
      } catch (err) {
        logger.error("upload image failed", { kind: "upload", err: String(err) });
        setUploading(false);
        return;
      }
      setUploading(false);
    }
    if (!finalImagePath) {
      finalImagePath = "assets/default_image.jpg";
    }

    if (isEdit && product) {
      await onSubmit({
        mode: "update",
        id: product.product_id,
        values: {
          product_name: values.product_name,
          product_quantity: values.product_quantity,
          product_prices: values.product_prices,
          product_type: values.product_type,
          image_path: finalImagePath,
        },
      });
    } else {
      await onSubmit({
        mode: "create",
        values: {
          ...values,
          image_path: finalImagePath,
          created_by: currentUserId,
        },
      });
    }
    onOpenChange(false);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit product" : "Create product"}</DialogTitle>
          <DialogDescription>
            Validated with Zod. Image upload is staged on the server before the
            product record is persisted.
          </DialogDescription>
        </DialogHeader>

        <form
          onSubmit={handleSubmit(handleFormSubmit)}
          className="space-y-4"
          id="product-form"
        >
          {!isEdit ? (
            <div className="space-y-2">
              <Label htmlFor="product_id">SKU</Label>
              <Input
                id="product_id"
                placeholder="SKU10001"
                {...register("product_id")}
              />
              {errors.product_id ? (
                <p className="text-xs text-destructive">
                  {errors.product_id.message}
                </p>
              ) : null}
            </div>
          ) : null}

          <div className="space-y-2">
            <Label htmlFor="product_name">Name</Label>
            <Input id="product_name" {...register("product_name")} />
            {errors.product_name ? (
              <p className="text-xs text-destructive">
                {errors.product_name.message}
              </p>
            ) : null}
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <div className="space-y-2">
              <Label htmlFor="product_quantity">Quantity</Label>
              <Input
                id="product_quantity"
                type="number"
                min={0}
                {...register("product_quantity", { valueAsNumber: true })}
              />
              {errors.product_quantity ? (
                <p className="text-xs text-destructive">
                  {errors.product_quantity.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="product_prices">Price</Label>
              <Input
                id="product_prices"
                placeholder="29.99"
                {...register("product_prices")}
              />
              {errors.product_prices ? (
                <p className="text-xs text-destructive">
                  {errors.product_prices.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="product_type">Type</Label>
              <Controller
                control={control}
                name="product_type"
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger id="product_type">
                      <SelectValue placeholder="Select type" />
                    </SelectTrigger>
                    <SelectContent>
                      {TYPE_OPTIONS.map((opt) => (
                        <SelectItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label>Image</Label>
            <ImageDropzone
              value={file}
              onChange={(f) => {
                setFile(f);
                if (f) setValue("image_path", "", { shouldValidate: false });
              }}
              initialPreviewUrl={
                imagePathValue ? resolveImageUrl(imagePathValue) : null
              }
              disabled={isSubmitting || uploading}
              error={errors.image_path?.message}
            />
            <input type="hidden" {...register("image_path")} />
          </div>
        </form>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSubmitting || uploading}
          >
            Cancel
          </Button>
          <Button
            form="product-form"
            type="submit"
            disabled={isSubmitting || uploading}
          >
            {uploading
              ? "Uploading…"
              : isSubmitting
              ? "Saving…"
              : isEdit
              ? "Save changes"
              : "Create product"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
