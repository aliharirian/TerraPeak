import React from "react";

const HomeDocs = () => {
  return (
    <section id="docs" className="py-20 px-4 sm:px-6 lg:px-8 bg-gray-900/50">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Documentation & Resources
          </h2>
          <p className="text-lg text-gray-400">
            Everything you need to get the most out of TerraPeak
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {/* GitHub Repo */}
          <a
            href="https://github.com/aliharirian/TerraPeak"
            target="_blank"
            rel="noopener noreferrer"
            className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-gray-900 rounded-lg flex items-center justify-center flex-shrink-0 group-hover:bg-green-500/10 transition-colors">
                <svg
                  className="w-6 h-6 text-gray-400 group-hover:text-green-500 transition-colors"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    fillRule="evenodd"
                    d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-green-500 transition-colors">
                  GitHub Repository
                </h3>
                <p className="text-gray-400 text-sm">
                  View source code, report issues, and contribute to the project
                </p>
              </div>
            </div>
          </a>

          {/* README */}
          <a
            href="https://github.com/aliharirian/TerraPeak#readme"
            target="_blank"
            rel="noopener noreferrer"
            className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-gray-900 rounded-lg flex items-center justify-center flex-shrink-0 group-hover:bg-green-500/10 transition-colors">
                <svg
                  className="w-6 h-6 text-gray-400 group-hover:text-green-500 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-green-500 transition-colors">
                  README Documentation
                </h3>
                <p className="text-gray-400 text-sm">
                  Complete setup guide, configuration options, and usage
                  examples
                </p>
              </div>
            </div>
          </a>

          {/* Issues */}
          <a
            href="https://github.com/aliharirian/TerraPeak/issues"
            target="_blank"
            rel="noopener noreferrer"
            className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-gray-900 rounded-lg flex items-center justify-center flex-shrink-0 group-hover:bg-green-500/10 transition-colors">
                <svg
                  className="w-6 h-6 text-gray-400 group-hover:text-green-500 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-green-500 transition-colors">
                  Report Issues
                </h3>
                <p className="text-gray-400 text-sm">
                  Found a bug or have a feature request? Let us know on GitHub
                </p>
              </div>
            </div>
          </a>

          {/* License */}
          <a
            href="https://github.com/aliharirian/TerraPeak/blob/main/LICENSE"
            target="_blank"
            rel="noopener noreferrer"
            className="bg-black border border-gray-800 rounded-xl p-6 hover:border-green-500/50 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 group"
          >
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 bg-gray-900 rounded-lg flex items-center justify-center flex-shrink-0 group-hover:bg-green-500/10 transition-colors">
                <svg
                  className="w-6 h-6 text-gray-400 group-hover:text-green-500 transition-colors"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                  />
                </svg>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-green-500 transition-colors">
                  Open Source License
                </h3>
                <p className="text-gray-400 text-sm">
                  Free to use, modify, and distribute under open source license
                </p>
              </div>
            </div>
          </a>
        </div>

        {/* Badges */}
        <div className="mt-12 flex flex-wrap items-center justify-center gap-4">
          <a
            href="https://github.com/aliharirian/TerraPeak"
            target="_blank"
            rel="noopener noreferrer"
            className="focus:outline-none focus:ring-2 focus:ring-green-500 rounded"
          >
            <img
              src="https://img.shields.io/github/stars/aliharirian/TerraPeak?style=social"
              alt="GitHub stars"
              className="h-5"
            />
          </a>
          <a
            href="https://github.com/aliharirian/TerraPeak/blob/main/LICENSE"
            target="_blank"
            rel="noopener noreferrer"
            className="focus:outline-none focus:ring-2 focus:ring-green-500 rounded"
          >
            <img
              src="https://img.shields.io/github/license/aliharirian/TerraPeak"
              alt="License"
              className="h-5"
            />
          </a>
          <a
            href="https://github.com/aliharirian/TerraPeak/issues"
            target="_blank"
            rel="noopener noreferrer"
            className="focus:outline-none focus:ring-2 focus:ring-green-500 rounded"
          >
            <img
              src="https://img.shields.io/github/issues/aliharirian/TerraPeak"
              alt="GitHub issues"
              className="h-5"
            />
          </a>
        </div>
      </div>
    </section>
  );
};

export default HomeDocs;
