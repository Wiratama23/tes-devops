import { DashboardStats } from "@/components/admin/DashboardStats";

export const metadata = {
  title: "Dashboard",
};

// Forced dynamic so the dashboard always reflects the latest data. CSR is
// powered by TanStack Query inside <DashboardStats />.
export const dynamic = "force-dynamic";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <header className="space-y-1">
        <p className="text-sm uppercase tracking-wide text-muted-foreground">
          Overview · CSR
        </p>
        <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Live counts pulled from the Go API via TanStack Query.
        </p>
      </header>

      <DashboardStats />
    </div>
  );
}
