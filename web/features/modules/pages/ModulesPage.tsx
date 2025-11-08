"use client";

import { useState } from "react";
import BrowseHeader from "../components/ModuleHeader";
import BrowseFooter from "../components/ModuleFooter";
import ModuleHero from "../components/ModuleHero";
import ModuleDetailModal from "../components/ModuleDetail/ModuleDetailModal";
import { TerraformModule } from "../types/moduleType";
import ModuleList from "../components/ModuleList/ModuleList";

// This would normally be in a separate metadata export, but since this is a client component,
// we'll set it via Head or in the layout

// Sample module data

export default function ModulesPage() {
  const [heroSearchTerm, setHeroSearchTerm] = useState("");
  const [selectedModule, setSelectedModule] = useState<TerraformModule | null>(
    null
  );
  const [searchTerm, setSearchTerm] = useState("");

  const handleHeroSearch = () => {
    setSearchTerm(heroSearchTerm);
    document
      .getElementById("module-list")
      ?.scrollIntoView({ behavior: "smooth" });
  };

  return (
    <div className="bg-gray-50 min-h-screen">
      {/* Header */}
      <BrowseHeader />

      {/* Hero Section */}
      <ModuleHero
        handleHeroSearch={handleHeroSearch}
        heroSearchTerm={heroSearchTerm}
        setHeroSearchTerm={(val) => setHeroSearchTerm(val)}
      />

      {/* Main Content */}
      <ModuleList
        searchTerm={searchTerm}
        setSearchTerm={(val) => setSearchTerm(val)}
        setSelectedModule={(module) => setSelectedModule(module)}
      />

      {/* Module Detail Modal */}
      {selectedModule && (
        <ModuleDetailModal
          selectedModule={selectedModule}
          setSelectedModule={(module) => setSelectedModule(module)}
        />
      )}

      {/* Footer */}
      <BrowseFooter />
    </div>
  );
}
