import "@testing-library/jest-dom/vitest";
import { afterAll, afterEach, beforeAll, vi } from "vitest";

import { server } from "./msw/server";

vi.stubEnv("NEXT_PUBLIC_API_BASE_URL", "http://api.test/api");
vi.stubEnv("INTERNAL_API_URL", "http://api.test/api");
vi.stubEnv("NEXT_PUBLIC_DEFAULT_IMAGE", "/assets/default_image.jpg");

beforeAll(() => server.listen({ onUnhandledRequest: "error" }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

if (typeof window !== "undefined") {
  if (!("matchMedia" in window)) {
    Object.defineProperty(window, "matchMedia", {
      writable: true,
      value: (query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: () => {},
        removeListener: () => {},
        addEventListener: () => {},
        removeEventListener: () => {},
        dispatchEvent: () => false,
      }),
    });
  }
  if (!("scrollTo" in window)) {
    (window as unknown as { scrollTo: () => void }).scrollTo = () => {};
  }
}
