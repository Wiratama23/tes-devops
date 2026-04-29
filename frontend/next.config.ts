import type { NextConfig } from "next";
import path from "node:path";

const nextConfig: NextConfig = {
  turbopack: {
    // Pin the project root to this directory so the workspace root warning
    // doesn't trigger when there's an unrelated lockfile higher up the tree.
    root: path.resolve(import.meta.dirname),
  },
  compiler: {
    // Strip every console.* call from the production bundle except
    // console.error — those represent genuine runtime failures that we still
    // want surfaced for the central logger to capture. Dev/test builds keep
    // all console output.
    removeConsole:
      process.env.NODE_ENV === "production" ? { exclude: ["error"] } : false,
  },
  images: {
    formats: ["image/avif", "image/webp"],
    remotePatterns: [
      {
        protocol: "http",
        hostname: "localhost",
      },
      {
        protocol: "http",
        hostname: "127.0.0.1",
      },
      {
        protocol: "https",
        hostname: "**",
      },
    ],
  },
  experimental: {
    optimizePackageImports: [
      "lucide-react",
      "@tanstack/react-query",
    ],
  },
  poweredByHeader: false,
};

export default nextConfig;
