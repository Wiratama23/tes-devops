import Link from "next/link";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { SmartImage } from "@/components/site/SmartImage";
import { cn } from "@/tools/utils";
import type { Product } from "@/types/api";

const TYPE_LABELS: Record<string, string> = {
  "10": "Drinks",
  "05": "Books",
  "20": "Electronics",
};

interface ProductCardProps {
  product: Product;
  className?: string;
  priority?: boolean;
}

export function ProductCard({ product, className, priority }: ProductCardProps) {
  const typeLabel = TYPE_LABELS[product.product_type] ?? `Type ${product.product_type}`;
  const price = Number(product.product_prices);
  const formatted = Number.isFinite(price)
    ? price.toLocaleString(undefined, { style: "currency", currency: "USD" })
    : `$${product.product_prices}`;

  return (
    <Card
      className={cn(
        "group h-full overflow-hidden transition-all hover:-translate-y-0.5 hover:shadow-md",
        className
      )}
    >
      <Link
        href={`/products/${encodeURIComponent(product.product_id)}`}
        className="block focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        prefetch
      >
        <div className="relative aspect-square w-full overflow-hidden bg-muted">
          <SmartImage
            imagePath={product.image_path}
            alt={product.product_name}
            fill
            sizes="(min-width: 1024px) 25vw, (min-width: 640px) 33vw, 50vw"
            className="object-cover transition-transform duration-500 group-hover:scale-[1.03]"
            priority={priority}
          />
        </div>
        <CardContent className="space-y-2 pt-4">
          <div className="flex items-center justify-between gap-2">
            <Badge variant="secondary">{typeLabel}</Badge>
            <span className="text-xs text-muted-foreground">
              {product.product_quantity} in stock
            </span>
          </div>
          <h3 className="line-clamp-2 text-sm font-medium leading-snug">
            {product.product_name}
          </h3>
        </CardContent>
        <CardFooter className="pt-0 text-base font-semibold">{formatted}</CardFooter>
      </Link>
    </Card>
  );
}
