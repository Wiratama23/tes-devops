import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function ProductNotFound() {
  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col items-center gap-4 px-4 py-24 text-center sm:px-6">
      <p className="text-sm uppercase tracking-wide text-muted-foreground">
        404
      </p>
      <h1 className="text-3xl font-semibold tracking-tight">
        Product not found
      </h1>
      <p className="text-muted-foreground">
        The product you were looking for has been removed or never existed.
      </p>
      <Button asChild>
        <Link href="/products">Back to products</Link>
      </Button>
    </div>
  );
}
