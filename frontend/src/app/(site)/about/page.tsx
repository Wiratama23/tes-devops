import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

export const dynamic = "force-static";

export const metadata = {
  title: "About",
  description: "About Tesdevops and the people building it.",
};

const VALUES = [
  {
    title: "Stable defaults",
    body: "We pick boring, dependable tools (Go, Postgres, Next.js) and bend them only when the payoff is concrete.",
  },
  {
    title: "Show your work",
    body: "Every page tells you what it's doing — Server Component, ISR, or CSR — so debugging is never a guessing game.",
  },
  {
    title: "Ship in seconds",
    body: "Optimistic UI updates, view transitions, and prefetched links keep the app feeling instant even on slow networks.",
  },
];

export default function AboutPage() {
  return (
    <div className="mx-auto w-full max-w-3xl px-4 py-16 sm:px-6">
      <header className="space-y-4">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          About us
        </p>
        <h1 className="text-4xl font-semibold tracking-tight">
          Modern storefront engineering, kept simple.
        </h1>
        <p className="text-lg leading-relaxed text-muted-foreground">
          Tesdevops is a reference implementation of a production-grade storefront
          built on Go, Next.js 16, Postgres, and Bun. It's intentionally
          opinionated so the codebase stays easy to onboard onto.
        </p>
      </header>

      <Separator className="my-12" />

      <section className="grid gap-4 sm:grid-cols-2">
        {VALUES.map((v) => (
          <Card key={v.title}>
            <CardHeader>
              <CardTitle className="text-lg">{v.title}</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground">
              {v.body}
            </CardContent>
          </Card>
        ))}
      </section>

      <Separator className="my-12" />

      <section className="prose prose-sm max-w-none dark:prose-invert">
        <h2>Stack</h2>
        <ul>
          <li>Go 1.26 · chi · pgx · Coraza WAF · jwtauth</li>
          <li>Next.js 16 (App Router) · React 19 · TypeScript 7 (tsgo)</li>
          <li>Tailwind v4 · Shadcn UI · Framer Motion · Sonner</li>
          <li>TanStack Query · Zod · React Hook Form · Tiptap</li>
        </ul>
      </section>
    </div>
  );
}
