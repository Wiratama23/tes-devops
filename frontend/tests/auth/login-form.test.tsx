import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { SWRConfig } from "swr";

import { LoginForm } from "@/components/admin/LoginForm";

const replace = vi.fn();
const refresh = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace, refresh }),
  useSearchParams: () => new URLSearchParams(""),
}));

function renderForm() {
  return render(
    <SWRConfig value={{ provider: () => new Map() }}>
      <LoginForm />
    </SWRConfig>
  );
}

describe("LoginForm", () => {
  it("validates required fields client-side", async () => {
    renderForm();
    await userEvent.click(screen.getByRole("button", { name: /sign in/i }));
    expect(
      await screen.findByText(/Username is required/i)
    ).toBeInTheDocument();
  });

  it("redirects on successful admin login", async () => {
    renderForm();
    await userEvent.type(screen.getByLabelText(/username/i), "admin");
    await userEvent.type(screen.getByLabelText(/password/i), "secret123");
    await userEvent.click(screen.getByRole("button", { name: /sign in/i }));
    await waitFor(() => {
      expect(replace).toHaveBeenCalledWith("/admin/dashboard");
    });
  });

  it("shows an error on bad credentials", async () => {
    renderForm();
    await userEvent.type(screen.getByLabelText(/username/i), "admin");
    await userEvent.type(screen.getByLabelText(/password/i), "wrongpass");
    await userEvent.click(screen.getByRole("button", { name: /sign in/i }));
    expect(
      await screen.findByText(/Invalid username or password/i)
    ).toBeInTheDocument();
  });
});
