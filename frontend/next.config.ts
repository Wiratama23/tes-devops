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
  // Optimize images more aggressively
  images: {
    formats: ["image/avif", "image/webp"],
    qualities: [25, 50, 75, 85],
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
        protocol: "http",
        hostname: "nginx",
      },
    ],
    // Cache images for 1 year (immutable)
    minimumCacheTTL: 31536000,
    deviceSizes: [640, 750, 828, 1080, 1200],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
  },
  experimental: {
    optimizePackageImports: [
      "lucide-react",
      "@radix-ui/react-dialog",
      "@radix-ui/react-select",
      "@radix-ui/react-tabs",
      "@radix-ui/react-tooltip",
    ],
    // Optimize bundle size
    optimizeCss: true,
  },
  poweredByHeader: false,
};

export default nextConfig;
