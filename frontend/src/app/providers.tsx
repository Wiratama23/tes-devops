"use client";

import { ThemeProvider } from "next-themes";
import { useEffect, useState } from "react";
import { SWRConfig } from "swr";

import { Toaster } from "@/components/ui/sonner";
import { me } from "@/services/client/auth";
import { clientFetch } from "@/tools/client-api";

function AuthInitializer({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Check if user is already logged in via cookie when app loads
    const checkAuth = async () => {
      try {
        await me();
        // User is authenticated, cookie will be sent automatically on future requests
      } catch {
        // Not authenticated, that's ok - they'll be redirected to login if needed
      }
    };

    checkAuth();
  }, []);

  return <>{children}</>;
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [swrConfig] = useState(() => ({
    fetcher: (key: string) => clientFetch(key),
    revalidateOnFocus: false,
    shouldRetryOnError: false,
  }));

  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <SWRConfig value={swrConfig}>
        <AuthInitializer>
          {children}
          <Toaster />
        </AuthInitializer>
      </SWRConfig>
    </ThemeProvider>
  );
}
