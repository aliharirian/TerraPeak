import React from "react";
interface ModuleHeroProps {
  heroSearchTerm: string;
  handleHeroSearch: () => void;
  setHeroSearchTerm: (value: string) => void;
}
const ModuleHero = ({
  heroSearchTerm,
  handleHeroSearch,
  setHeroSearchTerm,
}: ModuleHeroProps) => {
  return (
    <section className="bg-gradient-to-br from-green-50 to-green-100 border-b border-green-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16">
        <div className="text-center max-w-3xl mx-auto">
          <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">
            Browse Terraform Modules
          </h1>
          <p className="text-lg text-gray-700 mb-8">
            Discover and share Terraform modules and providers for your
            infrastructure
          </p>

          {/* Hero Search */}
          <div className="relative max-w-2xl mx-auto">
            <input
              type="text"
              value={heroSearchTerm}
              onChange={(e) => setHeroSearchTerm(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleHeroSearch()}
              placeholder="Search modules, providers..."
              className="w-full px-4 py-3 pr-24 rounded-lg border-2 border-green-300 focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-400 text-gray-900"
              aria-label="Search modules and providers"
            />
            <button
              onClick={handleHeroSearch}
              className="absolute right-2 top-1/2 -translate-y-1/2 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors font-medium"
            >
              Search
            </button>
          </div>
        </div>
      </div>
    </section>
  );
};

export default ModuleHero;
