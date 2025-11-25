import { HeroSection } from "@/components/hero-section"
import { BrandStory } from "@/components/brand-story"
import { ProblemSolution } from "@/components/problem-solution"
import { FeaturesGrid } from "@/components/features-grid"
import { QuickStart } from "@/components/quick-start"
import { IntegrationSection } from "@/components/integration-section"
import { CTASection } from "@/components/cta-section"
import { Footer } from "@/components/footer"

export default function HomePage() {
  return (
    <main className="min-h-screen bg-background">
      <HeroSection />
      <BrandStory />
      <ProblemSolution />
      <FeaturesGrid />
      <QuickStart />
      <IntegrationSection />
      <CTASection />
      <Footer />
    </main>
  )
}
