"use client";

import {
  FileText,
  LayoutDashboard,
  LogOut,
  Package,
  ShoppingBag,
} from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { toast } from "sonner";

import { ThemeToggle } from "@/components/site/ThemeToggle";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { logger } from "@/tools/logger";
import { cn } from "@/tools/utils";
import { logout } from "@/services/auth";

const NAV = [
  { href: "/admin/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/products", label: "Products", icon: Package },
  { href: "/admin/articles", label: "Articles", icon: FileText },
];

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();

  const logoutMutation = useMutation({
    mutationFn: logout,
    onSuccess: () => {
      toast.success("Signed out");
      router.replace("/admin/login");
      router.refresh();
    },
    onError: (err) => {
      logger.error("logout failed", { kind: "logout", err: String(err) });
      toast.error("Couldn't sign out — try again.");
    },
  });

  return (
    <div className="flex min-h-screen flex-col md:flex-row">
      <aside className="border-b border-border bg-card md:w-64 md:border-b-0 md:border-r">
        <div className="flex h-full flex-col gap-4 p-4">
          <Link
            href="/"
            className="flex items-center gap-2 px-2 py-2 text-sm font-semibold"
          >
            <ShoppingBag className="h-5 w-5 text-primary" />
            <span>Tesdevops Admin</span>
          </Link>
          <Separator />
          <nav className="flex flex-row gap-1 overflow-x-auto md:flex-col">
            {NAV.map((item) => {
              const active = pathname?.startsWith(item.href);
              const Icon = item.icon;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-2 rounded-md px-3 py-2 text-sm transition-colors",
                    active
                      ? "bg-secondary font-medium text-secondary-foreground"
                      : "text-muted-foreground hover:bg-muted hover:text-foreground"
                  )}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Link>
              );
            })}
          </nav>
          <div className="mt-auto flex items-center justify-between gap-2 pt-4">
            <ThemeToggle />
            <Button
              variant="ghost"
              size="sm"
              onClick={() => logoutMutation.mutate()}
              disabled={logoutMutation.isPending}
            >
              <LogOut className="mr-2 h-4 w-4" />
              Sign out
            </Button>
          </div>
        </div>
      </aside>
      <div className="flex-1 bg-background">
        <div className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6">{children}</div>
      </div>
    </div>
  );
}
