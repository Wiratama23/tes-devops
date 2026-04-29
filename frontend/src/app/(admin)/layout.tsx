import type { ReactNode } from "react";

// The admin route group has no chrome of its own — the AdminShell + login
// page handle their own UI. The middleware in src/middleware.ts already
// guards every /admin/* route except /admin/login.
export default function AdminGroupLayout({ children }: { children: ReactNode }) {
  return <>{children}</>;
}
