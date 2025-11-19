"use client";

import Footer from "@/components/layout/Footer/Footer";
import Header from "@/components/layout/Header/Header";
import HomeArchitecture from "@/features/home/components/HomeArchitecture";
import HomeConfiguration from "@/features/home/components/HomeConfiguration";
import HomeDocs from "@/features/home/components/HomeDocs";
import HomeFeatures from "@/features/home/components/HomeFeatures";
import HomeHero from "@/features/home/components/HomeHero";
import HomeQuickstart from "@/features/home/components/HomeQuickstart";
import { useState, useEffect } from "react";

export default function HomePage() {
  useEffect(() => {
    // Smooth scroll for anchor links
    const handleAnchorClick = (e: MouseEvent) => {
      const target = e.target as HTMLAnchorElement;
      if (target.tagName === "A" && target.hash) {
        const href = target.getAttribute("href");
        if (href?.startsWith("#")) {
          e.preventDefault();
          const element = document.querySelector(href);
          if (element) {
            const headerOffset = 80;
            const elementPosition = element.getBoundingClientRect().top;
            const offsetPosition =
              elementPosition + window.pageYOffset - headerOffset;
            window.scrollTo({
              top: offsetPosition,
              behavior: "smooth",
            });
          }
        }
      }
    };

    document.addEventListener("click", handleAnchorClick);
    return () => document.removeEventListener("click", handleAnchorClick);
  }, []);

  return (
    <div className="bg-black text-gray-100 font-sans antialiased">
      {/* Header / Navigation */}
      <Header />

      <main>
        {/* Hero Section */}
        <HomeHero />

        {/* Features Section */}
        <HomeFeatures />

        {/* Quickstart Section */}
        <HomeQuickstart />

        {/* Configuration Section */}
        <HomeConfiguration />

        {/* Architecture Section */}
        <HomeArchitecture />

        {/* Docs & Resources Section */}
        <HomeDocs />
      </main>

      {/* Footer */}
      <Footer />
    </div>
  );
}
