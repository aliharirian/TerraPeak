/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  // Enable standalone output for optimized Docker builds
  // This creates a minimal server.js and dependencies in .next/standalone
  output: 'standalone',
}

export default nextConfig
