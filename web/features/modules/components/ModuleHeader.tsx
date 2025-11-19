import { Menu, X } from "lucide-react";
import React, { useState } from "react";

const BrowseHeader = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <a href="/" className="flex items-center gap-2">
              <span className="w-8 h-8 bg-green-600 rounded-full flex items-center justify-center text-white font-bold text-sm">
                T
              </span>
              <span className="text-xl font-bold text-gray-900">
                terra<span className="text-green-600">peak</span>
              </span>
            </a>
          </div>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-6">
            <a
              href="/"
              className="text-gray-700 hover:text-green-600 transition-colors font-medium"
            >
              Home
            </a>
            <a href="/browse" className="text-green-600 font-medium">
              Browse
            </a>
            <a
              href="/#features"
              className="text-gray-700 hover:text-green-600 transition-colors font-medium"
            >
              Features
            </a>
            <a
              href="/#docs"
              className="text-gray-700 hover:text-green-600 transition-colors font-medium"
            >
              Documentation
            </a>
          </div>

          {/* Mobile Menu Button */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-md text-gray-700 hover:bg-gray-100"
            aria-label="Toggle menu"
            aria-expanded={mobileMenuOpen}
          >
            {mobileMenuOpen ? (
              <X className="w-6 h-6" />
            ) : (
              <Menu className="w-6 h-6" />
            )}
          </button>
        </div>

        {/* Mobile Nav */}
        {mobileMenuOpen && (
          <div className="md:hidden pb-4">
            <div className="flex flex-col gap-2">
              <a
                href="/"
                className="px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 hover:text-green-600 transition-colors"
              >
                Home
              </a>
              <a
                href="/browse"
                className="px-3 py-2 rounded-md text-green-600 bg-green-50 transition-colors"
              >
                Browse
              </a>
              <a
                href="/#features"
                className="px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 hover:text-green-600 transition-colors"
              >
                Features
              </a>
              <a
                href="/#docs"
                className="px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 hover:text-green-600 transition-colors"
              >
                Documentation
              </a>
            </div>
          </div>
        )}
      </nav>
    </header>
  );
};

export default BrowseHeader;
