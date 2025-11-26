# =============================================================================
# TerraPeak Landing - Production-Ready Multi-Stage Dockerfile
# =============================================================================
# Optimized Next.js deployment with:
# - Multi-stage build for minimal image size (~150MB)
# - Advanced layer caching for faster rebuilds
# - Security-hardened with non-root user
# - Optimized for pnpm with proper caching
# - Environment variable support via ARGs
# =============================================================================

# Build arguments for flexibility
ARG NODE_VERSION=20
ARG PNPM_VERSION=9

# -----------------------------------------------------------------------------
# Stage 1: Dependencies
# Install only production dependencies for better caching
# -----------------------------------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS deps

ARG PNPM_VERSION

# Install system dependencies
RUN apk add --no-cache libc6-compat

# Install pnpm globally
RUN corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate

WORKDIR /app

# Copy dependency manifests
COPY package.json pnpm-lock.yaml ./

# Install dependencies with efficient caching
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile --prefer-offline

# -----------------------------------------------------------------------------
# Stage 2: Builder
# Build the Next.js application with all optimizations
# -----------------------------------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS builder

ARG PNPM_VERSION

# Install pnpm
RUN corepack enable && corepack prepare pnpm@${PNPM_VERSION} --activate

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy application source (respects .dockerignore)
COPY . .
COPY ./.env.example ./.env

# Build-time environment variables
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Build the application
# Creates optimized standalone output in .next/standalone
RUN pnpm build

# -----------------------------------------------------------------------------
# Stage 3: Runner (Production)
# Minimal runtime image with only necessary files
# -----------------------------------------------------------------------------
FROM node:${NODE_VERSION}-alpine AS runner

# Metadata labels
LABEL maintainer="TerraPeak Team"
LABEL org.opencontainers.image.title="TerraPeak Landing"
LABEL org.opencontainers.image.description="TerraPeak marketing and documentation landing page"
LABEL org.opencontainers.image.vendor="TerraPeak"

WORKDIR /app

# Production environment variables
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Install minimal runtime dependencies
RUN apk add --no-cache curl

# Create non-root user for security
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs && \
    mkdir -p .next && \
    chown -R nextjs:nodejs .next

# Copy public assets (static files)
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Copy standalone output (includes all necessary runtime files)
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Switch to non-root user
USER nextjs

# Expose application port
EXPOSE 3000

# Health check using the dedicated endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:3000/api/health || exit 1

# Start the Next.js production server
CMD ["node", "server.js"]
