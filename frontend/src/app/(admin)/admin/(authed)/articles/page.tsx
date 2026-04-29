import { ArticlesAdmin } from "@/components/admin/ArticlesAdmin";

export const dynamic = "force-dynamic";

export const metadata = {
  title: "Articles · Admin",
};

export default function AdminArticlesPage() {
  return (
    <div className="space-y-8">
      <header className="space-y-1">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Articles · CSR
        </p>
        <h1 className="text-3xl font-semibold tracking-tight">Articles</h1>
        <p className="text-muted-foreground">
          Compose articles in the rich-text editor. Deletes are optimistic.
        </p>
      </header>

      <ArticlesAdmin />
    </div>
  );
}
