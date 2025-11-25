import { Github, Mountain } from "lucide-react"
import Link from "next/link"

export function Footer() {
  return (
    <footer className="border-t border-border bg-background">
      <div className="max-w-7xl mx-auto px-6 sm:px-8 md:px-10 py-8 sm:py-10 md:py-12">
        {/* Main footer content */}
        <div className="mb-8 sm:mb-10">
          {/* Mobile: Brand full width on top, Resources & Community in 2 columns below */}
          {/* Desktop (md+): All 3 sections in a single row */}
          <div className="grid grid-cols-1 md:grid-cols-[1.5fr_1fr_1fr] gap-6 sm:gap-8 md:gap-10 lg:gap-12">

            {/* Brand Section */}
            <div className="space-y-2 sm:space-y-3">
              <div className="flex items-center gap-2.5">
                <Mountain className="w-8 h-8 sm:w-9 sm:h-9 md:w-10 md:h-10 text-primary" />
                <span className="text-xl sm:text-2xl md:text-2xl font-bold">TerraPeak</span>
              </div>
              <p className="text-xs sm:text-sm text-muted-foreground leading-relaxed">
                High-performance caching proxy for Terraform Registry. Built by Peaker.
              </p>
            </div>

            {/* Resources & Community - Grid container for mobile 2-col layout */}
            <div className="grid grid-cols-2 md:contents gap-5 sm:gap-6 md:gap-0">

              {/* Resources Section */}
              <div>
                <h4 className="font-semibold mb-2.5 sm:mb-3 text-sm md:text-base">Resources</h4>
                <ul className="space-y-1.5 sm:space-y-2 text-xs sm:text-sm text-muted-foreground">
                  <li>
                    <Link
                      href="https://github.com/aliharirian/TerraPeak"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      GitHub Repository
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="https://github.com/aliharirian/TerraPeak/blob/main/LICENSE"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      License (Apache 2.0)
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="/demo"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      Demo
                    </Link>
                  </li>
                </ul>
              </div>

              {/* Community Section */}
              <div>
                <h4 className="font-semibold mb-2.5 sm:mb-3 text-sm md:text-base">Community</h4>
                <ul className="space-y-1.5 sm:space-y-2 text-xs sm:text-sm text-muted-foreground">
                  <li>
                    <Link
                      href="https://github.com/aliharirian/TerraPeak/blob/main/CONTRIBUTING.md"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      Contributing
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="https://github.com/aliharirian/TerraPeak/issues"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      Report Issues
                    </Link>
                  </li>
                  <li>
                    <Link
                      href="https://github.com/aliharirian/TerraPeak/discussions"
                      className="hover:text-primary transition-colors inline-block"
                    >
                      Discussions
                    </Link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom bar - Copyright and social */}
        <div className="pt-6 sm:pt-7 md:pt-8 border-t border-border flex flex-row items-center justify-center md:justify-between gap-3 sm:gap-4">
          <p className="text-xs sm:text-sm text-muted-foreground text-center md:text-left">
            Made with ❤️ by the Peaker Team
          </p>
          <Link
            href="https://github.com/aliharirian/TerraPeak"
            className="text-muted-foreground hover:text-primary transition-colors"
            aria-label="GitHub Repository"
          >
            <Github className="w-5 h-5 sm:w-5 sm:h-5" />
          </Link>
        </div>
      </div>
    </footer>
  )
}
