import React, { useState } from "react";

const Header = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-black/80 backdrop-blur-md border-b border-gray-800">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
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
            <span className="text-xl font-bold text-white">terrapeak</span>
          </div>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-8">
            <a
              href="#home"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              Home
            </a>
            <a
              href="/browse"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              Browse
            </a>
            <a
              href="#quickstart"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              Quickstart
            </a>
            <a
              href="#features"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              Features
            </a>
            <a
              href="#docs"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              Docs
            </a>
            <a
              href="https://github.com/aliharirian/TerraPeak"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-300 hover:text-white transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black rounded px-2 py-1"
            >
              GitHub
            </a>
          </div>

          {/* Star Button */}
          <a
            href="https://github.com/aliharirian/TerraPeak"
            target="_blank"
            rel="noopener noreferrer"
            className="hidden sm:flex items-center gap-2 px-4 py-2 bg-gray-900 hover:bg-gray-800 text-white rounded-lg border border-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black"
          >
            <svg
              className="w-4 h-4"
              fill="currentColor"
              viewBox="0 0 20 20"
              aria-hidden="true"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
            <span className="text-sm font-medium">Star on GitHub</span>
          </a>

          {/* Mobile Menu Button */}
          <button
            id="mobile-menu-btn"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 text-gray-300 hover:text-white focus:outline-none focus:ring-2 focus:ring-green-500 rounded"
            aria-label="Toggle menu"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div id="mobile-menu" className="md:hidden pb-4">
            <div className="flex flex-col gap-2">
              <a
                href="#home"
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Home
              </a>
              <a
                href="/browse"
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Browse
              </a>
              <a
                href="#quickstart"
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Quickstart
              </a>
              <a
                href="#features"
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Features
              </a>
              <a
                href="#docs"
                onClick={() => setMobileMenuOpen(false)}
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Docs
              </a>
              <a
                href="https://github.com/aliharirian/TerraPeak"
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-300 hover:text-white transition-colors px-2 py-2 rounded focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                GitHub
              </a>
            </div>
          </div>
        )}
      </nav>
    </header>
  );
};

export default Header;
