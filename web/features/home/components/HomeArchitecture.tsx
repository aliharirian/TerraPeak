import React from "react";

const HomeArchitecture = () => {
  return (
    <section className="py-20 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            How It Works
          </h2>
          <p className="text-lg text-gray-400">
            Simple, efficient architecture for maximum performance
          </p>
        </div>

        <div className="bg-black border border-gray-800 rounded-xl p-8">
          {/* Architecture Diagram */}
          <div className="flex flex-col md:flex-row items-center justify-center gap-8 mb-8">
            {/* Terraform Client */}
            <div className="flex flex-col items-center">
              <div className="w-20 h-20 bg-purple-500/10 border-2 border-purple-500 rounded-lg flex items-center justify-center mb-3">
                <svg
                  className="w-10 h-10 text-purple-500"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path d="M12 2L2 7v10c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V7l-10-5z" />
                </svg>
              </div>
              <span className="text-sm font-semibold text-white">
                Terraform Client
              </span>
            </div>

            {/* Arrow */}
            <svg
              className="w-8 h-8 text-green-500 rotate-90 md:rotate-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 7l5 5m0 0l-5 5m5-5H6"
              />
            </svg>

            {/* TerraPeak */}
            <div className="flex flex-col items-center">
              <div className="w-20 h-20 bg-green-500/10 border-2 border-green-500 rounded-lg flex items-center justify-center mb-3">
                <svg
                  className="w-10 h-10 text-green-500"
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
              <span className="text-sm font-semibold text-white">
                TerraPeak
              </span>
              <span className="text-xs text-gray-500">(Cache Layer)</span>
            </div>

            {/* Arrow */}
            <svg
              className="w-8 h-8 text-green-500 rotate-90 md:rotate-0"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 7l5 5m0 0l-5 5m5-5H6"
              />
            </svg>

            {/* Storage / Registry */}
            <div className="flex flex-col items-center">
              <div className="w-20 h-20 bg-blue-500/10 border-2 border-blue-500 rounded-lg flex items-center justify-center mb-3">
                <svg
                  className="w-10 h-10 text-blue-500"
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
              <span className="text-sm font-semibold text-white">
                MinIO / Registry
              </span>
              <span className="text-xs text-gray-500">(Backend Storage)</span>
            </div>
          </div>

          {/* Explanation */}
          <div className="border-t border-gray-800 pt-6">
            <h3 className="text-lg font-semibold text-white mb-3">
              Request Flow
            </h3>
            <ol className="space-y-3 text-gray-400">
              <li className="flex gap-3">
                <span className="flex-shrink-0 w-6 h-6 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center text-xs font-bold">
                  1
                </span>
                <span>Terraform client requests a provider from TerraPeak</span>
              </li>
              <li className="flex gap-3">
                <span className="flex-shrink-0 w-6 h-6 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center text-xs font-bold">
                  2
                </span>
                <span>
                  TerraPeak checks its cache (MinIO or local storage) for the
                  provider
                </span>
              </li>
              <li className="flex gap-3">
                <span className="flex-shrink-0 w-6 h-6 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center text-xs font-bold">
                  3
                </span>
                <span>
                  If cached, TerraPeak serves the provider immediately (fast
                  path)
                </span>
              </li>
              <li className="flex gap-3">
                <span className="flex-shrink-0 w-6 h-6 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center text-xs font-bold">
                  4
                </span>
                <span>
                  If not cached, TerraPeak fetches from upstream registry,
                  caches it, and serves to client
                </span>
              </li>
              <li className="flex gap-3">
                <span className="flex-shrink-0 w-6 h-6 bg-green-500/10 text-green-500 rounded-full flex items-center justify-center text-xs font-bold">
                  5
                </span>
                <span>
                  Subsequent requests for the same provider are served from
                  cache, dramatically reducing download times
                </span>
              </li>
            </ol>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HomeArchitecture;
