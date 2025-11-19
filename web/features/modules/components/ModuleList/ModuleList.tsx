import React, { useEffect, useState } from "react";
import { TerraformModule } from "../../types/moduleType";
import { ChevronRight, Star, X } from "lucide-react";
import { MODULES_DATA } from "../../data";

interface ModuleListProps {
  searchTerm: string;
  setSearchTerm: (val: string) => void;
  setSelectedModule: (module: TerraformModule | null) => void;
}

const ModuleList: React.FC<ModuleListProps> = ({
  searchTerm,
  setSearchTerm,
  setSelectedModule,
}) => {
  const [filteredModules, setFilteredModules] =
    useState<TerraformModule[]>(MODULES_DATA);
  const [displayedCount, setDisplayedCount] = useState(6);
  const [selectedTags, setSelectedTags] = useState<Set<string>>(new Set());
  const [selectedProviders, setSelectedProviders] = useState<Set<string>>(
    new Set()
  );
  const [currentSort, setCurrentSort] = useState("stars");

  useEffect(() => {
    const filtered = MODULES_DATA.filter((module) => {
      const matchesSearch =
        !searchTerm ||
        module.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        module.namespace.toLowerCase().includes(searchTerm.toLowerCase()) ||
        module.description.toLowerCase().includes(searchTerm.toLowerCase());

      const matchesTags =
        selectedTags.size === 0 ||
        module.tags.some((tag) => selectedTags.has(tag));

      const matchesProvider =
        selectedProviders.size === 0 || selectedProviders.has(module.provider);

      return matchesSearch && matchesTags && matchesProvider;
    });

    filtered.sort((a, b) => {
      if (currentSort === "stars") return b.stars - a.stars;
      if (currentSort === "recent")
        return (
          new Date(b.versions[0].date).getTime() -
          new Date(a.versions[0].date).getTime()
        );
      if (currentSort === "name") return a.name.localeCompare(b.name);
      return 0;
    });

    setFilteredModules(filtered);
    setDisplayedCount(6);
  }, [searchTerm, selectedTags, selectedProviders, currentSort]);

  const toggleTag = (tag: string) => {
    const newTags = new Set(selectedTags);
    newTags.has(tag) ? newTags.delete(tag) : newTags.add(tag);
    setSelectedTags(newTags);
  };

  const toggleProvider = (provider: string) => {
    const newProviders = new Set(selectedProviders);
    newProviders.has(provider)
      ? newProviders.delete(provider)
      : newProviders.add(provider);
    setSelectedProviders(newProviders);
  };

  const clearFilters = () => {
    setSelectedTags(new Set());
    setSelectedProviders(new Set());
  };
  const allTags = Array.from(
    new Set(MODULES_DATA.flatMap((m) => m.tags.map((t) => t.toLowerCase())))
  ).sort();

  const allProviders = Array.from(
    new Set(MODULES_DATA.map((m) => m.provider))
  ).sort();
  return (
    <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex flex-col lg:flex-row gap-8">
        {/* Sidebar */}
        <aside className="lg:w-64 flex-shrink-0">
          <div className="bg-white rounded-lg border border-gray-200 p-4 sticky top-20">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold text-gray-900">Filters</h2>
              <button
                onClick={clearFilters}
                className="text-sm text-green-600 hover:text-green-700 font-medium"
              >
                Clear
              </button>
            </div>

            {/* Tags */}
            <div className="mb-6">
              <h3 className="text-sm font-semibold text-gray-700 mb-3">Tags</h3>
              <div className="flex flex-wrap gap-2">
                {allTags.map((tag) => (
                  <button
                    key={tag}
                    onClick={() => toggleTag(tag)}
                    className={`px-3 py-1 rounded-full border text-sm transition-colors ${
                      selectedTags.has(tag)
                        ? "border-green-500 bg-green-50 text-green-700"
                        : "border-gray-300 text-gray-700 hover:border-green-500 hover:bg-green-50 hover:text-green-700"
                    }`}
                  >
                    {tag}
                  </button>
                ))}
              </div>
            </div>

            {/* Providers */}
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-3">
                Providers
              </h3>
              <div className="space-y-2">
                {allProviders.map((provider) => (
                  <label
                    key={provider}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={selectedProviders.has(provider)}
                      onChange={() => toggleProvider(provider)}
                      className="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
                    />
                    <span className="text-sm text-gray-700">{provider}</span>
                  </label>
                ))}
              </div>
            </div>
          </div>
        </aside>

        {/* Module Cards */}
        <div className="flex-1">
          {/* Search + Sort */}
          <div className="bg-white rounded-lg border border-gray-200 p-4 mb-6">
            <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
              <div className="relative flex-1 w-full">
                <input
                  type="text"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  placeholder="Filter modules..."
                  className="w-full px-4 py-2 pr-10 rounded-md border border-gray-300 focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-400"
                />
                {searchTerm && (
                  <button
                    onClick={() => setSearchTerm("")}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  >
                    <X className="w-5 h-5" />
                  </button>
                )}
              </div>

              <div className="flex items-center gap-2">
                <label
                  htmlFor="sort-select"
                  className="text-sm text-gray-700 font-medium whitespace-nowrap"
                >
                  Sort by:
                </label>
                <select
                  id="sort-select"
                  value={currentSort}
                  onChange={(e) => setCurrentSort(e.target.value)}
                  className="px-3 py-2 rounded-md border border-gray-300 focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-400 text-sm"
                >
                  <option value="stars">Most Starred</option>
                  <option value="recent">Recently Updated</option>
                  <option value="name">Name (A-Z)</option>
                </select>
              </div>
            </div>
          </div>

          {/* Count */}
          <p className="text-sm text-gray-600 mb-4">
            Showing {Math.min(displayedCount, filteredModules.length)} of{" "}
            {filteredModules.length} modules
          </p>

          {/* Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {filteredModules.slice(0, displayedCount).map((module) => {
              const initials = module.namespace.substring(0, 2).toUpperCase();
              return (
                <article
                  key={module.id}
                  onClick={() => setSelectedModule(module)}
                  className="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-lg transition-shadow cursor-pointer"
                >
                  <div className="flex items-start gap-3 mb-3">
                    <div className="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center text-green-700 font-bold flex-shrink-0">
                      {initials}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-bold text-gray-900 truncate">
                        {module.namespace}/
                        <span className="text-green-600">{module.name}</span>
                      </h3>
                      <p className="text-sm text-gray-600 truncate">
                        {module.provider}
                      </p>
                    </div>
                  </div>

                  <p className="text-sm text-gray-700 mb-3 line-clamp-2">
                    {module.description}
                  </p>

                  <div className="flex flex-wrap gap-1 mb-3">
                    {module.tags.map((tag) => (
                      <span
                        key={tag}
                        className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded-full"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3 text-sm text-gray-600">
                      <span className="flex items-center gap-1">
                        <Star className="w-4 h-4 text-yellow-500 fill-yellow-500" />
                        {module.stars}
                      </span>
                      <span className="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full font-medium">
                        v{module.version}
                      </span>
                    </div>
                    <button className="px-3 py-1 text-sm text-green-600 hover:text-green-700 font-medium flex items-center gap-1">
                      View <ChevronRight className="w-4 h-4" />
                    </button>
                  </div>
                </article>
              );
            })}
          </div>

          {/* Load More */}
          {displayedCount < filteredModules.length && (
            <div className="text-center">
              <button
                onClick={() => setDisplayedCount(displayedCount + 6)}
                className="px-6 py-3 bg-white border-2 border-green-600 text-green-600 rounded-lg hover:bg-green-50 transition-colors font-medium"
              >
                Load More Modules
              </button>
            </div>
          )}
        </div>
      </div>
    </main>
  );
};

export default ModuleList;
