"use client";

import {
  FileText,
  LayoutDashboard,
  LogOut,
  Package,
  ShoppingBag,
} from "lucide-react";
import { useMutation, useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import { ThemeToggle } from "@/components/site/ThemeToggle";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { logger } from "@/tools/logger";
import { cn } from "@/tools/utils";
import { logout, me, refreshToken } from "@/services/auth";

const NAV = [
  { href: "/admin/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/admin/products", label: "Products", icon: Package },
  { href: "/admin/articles", label: "Articles", icon: FileText },
];

// Refresh token every 30 minutes if user is active, or after any user activity
const REFRESH_INTERVAL = 30 * 60 * 1000; // 30 minutes
const INACTIVITY_TIMEOUT = 55 * 60 * 1000; // 55 minutes (token expires in 60 minutes)
const ACTIVITY_THROTTLE = 60 * 1000; // Throttle activity-triggered refresh to at most once per minute

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const inactivityTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const refreshIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const lastActivityRefreshRef = useRef<number>(0);

  // Check if user is authenticated
  const { data: user, isLoading, isError } = useQuery({
    queryKey: ["auth", "me"],
    queryFn: () => me(),
    retry: false,
  });

  // Mutation to logout
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

  // Mutation to refresh token
  const refreshMutation = useMutation({
    mutationFn: refreshToken,
    onError: (err) => {
      logger.error("token refresh failed", { kind: "token_refresh", err: String(err) });
      // If refresh fails, redirect to login
      router.replace("/admin/login");
    },
  });

  // Setup token refresh and inactivity detection
  useEffect(() => {
    if (!user) return;

    // Refresh token periodically (every 30 minutes)
    refreshIntervalRef.current = setInterval(() => {
      refreshMutation.mutate();
    }, REFRESH_INTERVAL);

    // Handle user activity — only resets the inactivity timeout.
    // Triggers a token refresh at most once per ACTIVITY_THROTTLE interval
    // to avoid flooding the API on high-frequency events like scroll.
    const handleActivity = () => {
      // Reset the inactivity timeout on every event
      if (inactivityTimeoutRef.current) {
        clearTimeout(inactivityTimeoutRef.current);
      }
      inactivityTimeoutRef.current = setTimeout(() => {
        // If no activity for 55 minutes, logout
        logoutMutation.mutate();
      }, INACTIVITY_TIMEOUT);

      // Throttle the refresh call
      const now = Date.now();
      if (now - lastActivityRefreshRef.current >= ACTIVITY_THROTTLE) {
        lastActivityRefreshRef.current = now;
        refreshMutation.mutate();
      }
    };

    // Listen for user activity
    const events = ["mousedown", "keydown", "scroll", "touchstart", "click"];
    events.forEach((event) => {
      document.addEventListener(event, handleActivity);
    });

    // Initial inactivity timeout
    inactivityTimeoutRef.current = setTimeout(() => {
      logoutMutation.mutate();
    }, INACTIVITY_TIMEOUT);

    return () => {
      // Cleanup
      if (refreshIntervalRef.current) clearInterval(refreshIntervalRef.current);
      if (inactivityTimeoutRef.current) clearTimeout(inactivityTimeoutRef.current);
      events.forEach((event) => {
        document.removeEventListener(event, handleActivity);
      });
    };
  }, [user, logoutMutation, refreshMutation]);

  // Redirect to login if not authenticated
  useEffect(() => {
    if (!isLoading && isError) {
      router.replace("/admin/login");
    }
  }, [isLoading, isError, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-muted-foreground">Loading...</div>
      </div>
    );
  }

  if (isError || !user) {
    return null;
  }

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
