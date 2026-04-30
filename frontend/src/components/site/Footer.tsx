import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t border-border/60 bg-background">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-4 px-4 py-8 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <p>&copy; {new Date().getFullYear()} Tesdevops. Built with Next.js + Go.</p>
        <nav className="flex flex-wrap gap-4">
          <Link href="/about" className="hover:text-foreground transition-colors">
            About
          </Link>
          <Link href="/contact" className="hover:text-foreground transition-colors">
            Contact
          </Link>
          <Link href="/articles" className="hover:text-foreground transition-colors">
            Articles
          </Link>
          <Link href="/products" className="hover:text-foreground transition-colors">
            Products
          </Link>
        </nav>
      </div>
    </footer>
  );
}
