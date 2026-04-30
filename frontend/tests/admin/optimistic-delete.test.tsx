import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { delay, http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";

import { ArticlesAdmin } from "@/components/admin/ArticlesAdmin";
import { server } from "../msw/server";

function renderWithClient() {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={client}>
      <ArticlesAdmin />
    </QueryClientProvider>
  );
}

describe("ArticlesAdmin optimistic delete", () => {
  it("removes the row immediately and rolls back on failure", async () => {
    server.use(
      http.delete("http://api.test/api/articles/:id", async () => {
        // Slow the server down so we can observe the optimistic UI before the
        // 500 lands and triggers the rollback.
        await delay(150);
        return new HttpResponse("server fault", { status: 500 });
      })
    );

    renderWithClient();

    await screen.findByText("Getting Started");

    await userEvent.click(
      screen.getByRole("button", { name: /delete getting started/i })
    );
    await userEvent.click(
      await screen.findByRole("button", { name: /^delete$/i })
    );

    // Optimistic update: row should disappear right after confirming.
    await waitFor(() =>
      expect(screen.queryByText("Getting Started")).not.toBeInTheDocument()
    );

    // Rollback: row re-appears once the 500 response lands.
    await waitFor(
      () => expect(screen.getByText("Getting Started")).toBeInTheDocument(),
      { timeout: 2000 }
    );
  });
});
