import { ProductsAdmin } from "@/components/admin/ProductsAdmin";

export const dynamic = "force-dynamic";

export const metadata = {
  title: "Products · Admin",
};

export default function AdminProductsPage() {
  return (
    <div className="space-y-8">
      <header className="space-y-1">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Products · CSR
        </p>
        <h1 className="text-3xl font-semibold tracking-tight">Products</h1>
        <p className="text-muted-foreground">
          Create, edit, and delete products. Deletes are optimistic — the row
          disappears immediately and rolls back on failure.
        </p>
      </header>

      <ProductsAdmin />
    </div>
  );
}
