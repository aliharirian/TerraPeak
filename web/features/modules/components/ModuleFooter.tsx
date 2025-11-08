import React from "react";

const BrowseFooter = () => {
  return (
    <footer className="bg-white border-t border-gray-200 mt-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <div className="flex items-center gap-2 mb-4">
              <span className="w-8 h-8 bg-green-600 rounded-full flex items-center justify-center text-white font-bold text-sm">
                T
              </span>
              <span className="text-lg font-bold text-gray-900">
                terra<span className="text-green-600">peak</span>
              </span>
            </div>
            <p className="text-sm text-gray-600">
              Your trusted Terraform module registry
            </p>
          </div>

          <div>
            <h4 className="font-semibold text-gray-900 mb-3">Resources</h4>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="/#docs" className="text-gray-600 hover:text-green-600">
                  Documentation
                </a>
              </li>
              <li>
                <a
                  href="/#quickstart"
                  className="text-gray-600 hover:text-green-600"
                >
                  Quickstart
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/aliharirian/TerraPeak"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 hover:text-green-600"
                >
                  GitHub
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-gray-900 mb-3">Community</h4>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://github.com/aliharirian/TerraPeak"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 hover:text-green-600"
                >
                  GitHub
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/aliharirian/TerraPeak/issues"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 hover:text-green-600"
                >
                  Issues
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-gray-900 mb-3">Company</h4>
            <ul className="space-y-2 text-sm">
              <li>
                <a href="/" className="text-gray-600 hover:text-green-600">
                  About
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/aliharirian/TerraPeak/blob/main/LICENSE"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-gray-600 hover:text-green-600"
                >
                  License
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-200 mt-8 pt-8 text-center text-sm text-gray-600">
          <p>&copy; 2025 terrapeak. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
};

export default BrowseFooter;
