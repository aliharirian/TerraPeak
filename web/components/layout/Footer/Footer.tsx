import React from "react";

const Footer = () => {
  return (
    <footer className="py-12 px-4 sm:px-6 lg:px-8 border-t border-gray-800">
      <div className="max-w-7xl mx-auto">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-green-500 to-green-600 flex items-center justify-center">
              <svg
                className="w-5 h-5 text-white"
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
            <span className="text-lg font-bold text-white">terrapeak</span>
          </div>

          <div className="flex flex-col sm:flex-row items-center gap-4 text-sm text-gray-400">
            <span>© 2025 TerraPeak. Open Source Project.</span>
            <a
              href="https://github.com/aliharirian/TerraPeak/blob/main/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-green-500 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 rounded px-2 py-1"
            >
              License
            </a>
            <a
              href="https://github.com/aliharirian"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-green-500 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 rounded px-2 py-1"
            >
              Contributors
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
