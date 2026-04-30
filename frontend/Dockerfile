# Stage 1: Dependencies with oven/bun:canary-slim
FROM oven/bun:canary-slim AS deps

WORKDIR /app

# Copy package files
COPY package.json bun.lock* ./

# Install dependencies with frozen lockfile for reproducibility
RUN bun install --frozen-lockfile


# Stage 2: Builder with oven/bun:canary-slim
FROM oven/bun:canary-slim AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy all source files
COPY . .

# Disable Next.js telemetry during build
ENV NEXT_TELEMETRY_DISABLED=1

# Build the Next.js application
RUN bun run build


# Stage 3: Production runner with oven/bun:canary-slim
FROM oven/bun:canary-slim AS runner

WORKDIR /app

# Set production environment variables
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=4000

# Create non-root user for security
RUN groupadd --system --gid 1001 bun \
 && useradd  --system --uid 1001 --gid 1001 nextjs

# Copy public directory
COPY --from=builder /app/public ./public

# Copy Next.js build output with proper permissions
COPY --from=builder --chown=nextjs:bun /app/.next ./.next

# Copy node_modules with proper permissions
COPY --from=builder --chown=nextjs:bun /app/node_modules ./node_modules

# Copy package.json
COPY --from=builder /app/package.json ./package.json

# Switch to non-root user
USER nextjs

# Expose the application port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD bun run --silent -e "fetch('http://localhost:4000').catch(() => process.exit(1))"

# Start the application
CMD ["bun", "run", "start"]
