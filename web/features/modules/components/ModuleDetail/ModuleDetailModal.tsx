"use client";

import React, { useState } from "react";
import { X, Copy, Star } from "lucide-react";
import { TerraformModule } from "../../types/moduleType";

interface ModuleDetailModalProps {
  selectedModule: TerraformModule | null;
  setSelectedModule: (module: TerraformModule | null) => void;
}

const ModuleDetailModal: React.FC<ModuleDetailModalProps> = ({
  selectedModule,
  setSelectedModule,
}) => {
  const [copySuccess, setCopySuccess] = useState(false);

  if (!selectedModule) return null;

  const copyToClipboard = async () => {
    const text = `module "${selectedModule.name}" {
  source  = "${selectedModule.fullName}"
  version = "${selectedModule.version}"

  # Configuration options
}`;
    await navigator.clipboard.writeText(text);
    setCopySuccess(true);
    setTimeout(() => setCopySuccess(false), 2000);
  };

  return (
    <div
      onClick={() => setSelectedModule(null)}
      className="fixed inset-0 bg-black bg-opacity-50 z-50 overflow-y-auto"
    >
      <div className="min-h-screen px-4 py-8 flex items-center justify-center">
        <div
          onClick={(e) => e.stopPropagation()}
          className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto"
        >
          {/* Header */}
          <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
            <h2 className="text-2xl font-bold text-gray-900">
              {selectedModule.fullName}
            </h2>
            <button
              onClick={() => setSelectedModule(null)}
              className="p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
              aria-label="Close modal"
            >
              <X className="w-6 h-6" />
            </button>
          </div>

          {/* Content */}
          <div className="px-6 py-6">
            {/* Info */}
            <div className="mb-6">
              <div className="flex items-center gap-4 mb-4">
                <div className="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center text-green-700 font-bold text-xl">
                  {selectedModule.namespace.substring(0, 2).toUpperCase()}
                </div>
                <div>
                  <p className="text-sm text-gray-600">
                    {selectedModule.namespace}
                  </p>
                  <div className="flex items-center gap-3 mt-1">
                    <span className="flex items-center gap-1 text-sm text-gray-600">
                      <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                      {selectedModule.stars}
                    </span>
                    <span className="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full font-medium">
                      v{selectedModule.version}
                    </span>
                  </div>
                </div>
              </div>

              <p className="text-gray-700 leading-relaxed">
                {selectedModule.description}
              </p>

              <div className="flex flex-wrap gap-2 mt-4">
                {selectedModule.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-3 py-1 bg-gray-100 text-gray-700 text-sm rounded-full"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* Usage */}
            <div className="mb-6">
              <h3 className="text-lg font-bold text-gray-900 mb-3">Usage</h3>
              <div className="relative">
                <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto text-sm">
                  <code>{`module "${selectedModule.name}" {
  source  = "${selectedModule.fullName}"
  version = "${selectedModule.version}"

  # Configuration options
}`}</code>
                </pre>
                <button
                  onClick={copyToClipboard}
                  className="absolute top-3 right-3 px-3 py-1 bg-green-600 text-white text-sm rounded-md hover:bg-green-700 transition-colors flex items-center gap-1"
                >
                  <Copy className="w-4 h-4" />
                  {copySuccess ? "Copied!" : "Copy"}
                </button>
              </div>
            </div>

            {/* Versions */}
            <div>
              <h3 className="text-lg font-bold text-gray-900 mb-3">
                Available Versions
              </h3>
              <div className="space-y-2">
                {selectedModule.versions.map((v) => (
                  <div
                    key={v.version}
                    className="flex items-center justify-between p-3 bg-gray-50 rounded-md"
                  >
                    <div>
                      <span className="font-medium text-gray-900">
                        v{v.version}
                      </span>
                      <span className="text-sm text-gray-600 ml-3">
                        {v.date}
                      </span>
                    </div>
                    <button className="text-sm text-green-600 hover:text-green-700 font-medium">
                      Install
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Footer */}
          <div className="border-t border-gray-200 px-6 py-4 bg-gray-50">
            <button
              onClick={() => setSelectedModule(null)}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors font-medium"
            >
              Back to List
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ModuleDetailModal;
