"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { ThemeProvider } from "next-themes";
import { useEffect, useState } from "react";

import { Toaster } from "@/components/ui/sonner";
import { getQueryClient } from "@/tools/query-client";
import { me } from "@/services/auth";

function AuthInitializer({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    // Check if user is already logged in via cookie when app loads
    const checkAuth = async () => {
      try {
        await me();
        // User is authenticated, cookie will be sent automatically on future requests
      } catch (error) {
        // Not authenticated, that's ok - they'll be redirected to login if needed
      }
    };

    checkAuth();
  }, []);

  return <>{children}</>;
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [client] = useState(() => getQueryClient());

  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <QueryClientProvider client={client}>
        <AuthInitializer>
          {children}
          <Toaster />
          {process.env.NODE_ENV === "development" ? (
            <ReactQueryDevtools initialIsOpen={false} buttonPosition="bottom-left" />
          ) : null}
        </AuthInitializer>
      </QueryClientProvider>
    </ThemeProvider>
  );
}
