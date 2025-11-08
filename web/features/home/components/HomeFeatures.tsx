import React from "react";

const HomeFeatures = () => {
  return (
    <section
      id="features"
      className="py-20 px-4 sm:px-6 lg:px-8 bg-gray-900/50"
    >
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Why TerraPeak?
          </h2>
          <p className="text-lg text-gray-400 max-w-2xl mx-auto">
            Built for DevOps teams who need reliable, fast, and secure Terraform
            provider caching
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          {/* Feature 1 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              High-Performance Caching
            </h3>
            <p className="text-gray-400 leading-relaxed">
              Dramatically reduce provider download times and bandwidth usage
              with intelligent caching. Perfect for teams running frequent
              Terraform operations.
            </p>
          </div>

          {/* Feature 2 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              Intelligent Storage Backends
            </h3>
            <p className="text-gray-400 leading-relaxed">
              Choose between MinIO for distributed object storage or local file
              storage. Scale your caching infrastructure to match your needs.
            </p>
          </div>

          {/* Feature 3 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              Flexible Proxy Support
            </h3>
            <p className="text-gray-400 leading-relaxed">
              Works seamlessly in corporate environments with outbound client
              proxy and inbound server proxy modes. Supports HTTP, SOCKS5, and
              SOCKS4 protocols.
            </p>
          </div>

          {/* Feature 4 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              HTTPS Required
            </h3>
            <p className="text-gray-400 leading-relaxed">
              TerraPeak requires HTTPS with valid SSL certificates to ensure
              Terraform accepts provider downloads securely. Built with security
              as a priority.
            </p>
          </div>

          {/* Feature 5 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              Corporate-Friendly
            </h3>
            <p className="text-gray-400 leading-relaxed">
              Designed for enterprise environments with support for various
              proxy types and configurations. Integrate seamlessly with your
              existing infrastructure.
            </p>
          </div>

          {/* Feature 6 */}
          <div className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors">
            <div className="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-4">
              <svg
                className="w-6 h-6 text-green-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-3">
              Easy to Deploy
            </h3>
            <p className="text-gray-400 leading-relaxed">
              Get started in minutes with Docker Compose, Docker, or build from
              source. Simple configuration with YAML files for all settings.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HomeFeatures;
