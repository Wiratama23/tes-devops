import type { ReactNode } from "react";

import { AdminShell } from "@/components/admin/AdminShell";

// All admin routes nested under (authed) get the sidebar shell. Login lives
// outside this group so it can render its own minimal layout.
export default function AdminAuthedLayout({ children }: { children: ReactNode }) {
  return <AdminShell>{children}</AdminShell>;
}
