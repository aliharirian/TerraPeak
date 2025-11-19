# =============================================================================
# TerraPeak Web - Production-Ready Multi-Stage Dockerfile
# =============================================================================
# This Dockerfile follows Next.js best practices for production deployments
# - Multi-stage build for minimal final image size
# - Efficient layer caching for faster rebuilds
# - Security-hardened with non-root user
# - Optimized for pnpm package manager
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Dependencies
# Install dependencies in a separate stage for better caching
# -----------------------------------------------------------------------------
FROM node:20-alpine AS deps

# Install libc6-compat and pnpm for compatibility with some native dependencies
RUN apk add --no-cache libc6-compat && \
    npm install -g pnpm

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies with frozen lockfile for reproducible builds
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm config set store-dir /pnpm/store && \
    pnpm install --frozen-lockfile

# -----------------------------------------------------------------------------
# Stage 2: Builder
# Build the Next.js application
# -----------------------------------------------------------------------------
FROM node:20-alpine AS builder

# Install pnpm
RUN npm install -g pnpm

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy application source
COPY . .

# Set build-time environment variables
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Build the Next.js application
# This creates optimized production build in .next folder
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm config set store-dir /pnpm/store && \
    pnpm build

# -----------------------------------------------------------------------------
# Stage 3: Runner (Production)
# Run the application with minimal dependencies
# -----------------------------------------------------------------------------
FROM node:20-alpine AS runner

WORKDIR /app

# Set production environment
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Create a non-root user for security
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy necessary files from builder
# Copy public assets (if they exist)
COPY --from=builder /app/public ./public

# Copy package.json for runtime dependencies info
COPY --from=builder /app/package.json ./package.json

# Copy Next.js standalone output (includes minimal runtime)
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Health check to ensure container is ready
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD node -e "const http = require('http'); const options = { hostname: 'localhost', port: 3000, path: '/api/health', method: 'GET' }; const req = http.request(options, (res) => { process.exit(res.statusCode === 200 ? 0 : 1); }); req.on('error', () => process.exit(1)); req.end();"

# Start the Next.js server
# The server.js file is created by Next.js standalone build
CMD ["node", "server.js"]
