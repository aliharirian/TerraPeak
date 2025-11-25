import { Mountain, Plus } from "lucide-react"

export function BrandStory() {
  return (
    <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-background">
      <div className="max-w-6xl mx-auto">
        <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-center mb-8 sm:mb-12 md:mb-16 text-balance px-2">The Story Behind TerraPeak</h2>

        <div className="flex flex-row items-center justify-center gap-2 sm:gap-4 md:gap-6 lg:gap-8 xl:gap-12 overflow-x-auto px-2">
          {/* Terra */}
          <div className="flex flex-col items-center text-center space-y-1 sm:space-y-2 md:space-y-3 flex-shrink-0">
            <div className="w-12 h-12 sm:w-16 sm:h-16 md:w-20 md:h-20 lg:w-24 lg:h-24 rounded-lg sm:rounded-xl md:rounded-2xl bg-[#7B42BC]/20 flex items-center justify-center border-2 border-[#7B42BC]/40">
              <svg className="w-6 h-6 sm:w-8 sm:h-8 md:w-10 md:h-10 lg:w-12 lg:h-12" viewBox="0 0 24 24" fill="none">
                <rect x="3" y="3" width="8" height="8" fill="#7B42BC" />
                <rect x="13" y="3" width="8" height="8" fill="#7B42BC" opacity="0.7" />
                <rect x="3" y="13" width="8" height="8" fill="#7B42BC" opacity="0.7" />
                <rect x="13" y="13" width="8" height="8" fill="#7B42BC" />
              </svg>
            </div>
            <div>
              <div className="text-sm sm:text-base md:text-lg lg:text-xl xl:text-2xl font-bold text-[#7B42BC]">Terra</div>
              <div className="text-[10px] sm:text-xs md:text-sm text-muted-foreground hidden sm:block">from Terraform</div>
            </div>
          </div>

          {/* Plus */}
          <div className="flex items-center justify-center flex-shrink-0 -mb-0 sm:-mb-2 md:-mb-0">
            <Plus className="w-5 h-5 sm:w-6 sm:h-6 md:w-8 md:h-8 lg:w-10 lg:h-10 xl:w-12 xl:h-12 text-muted-foreground" />
          </div>

          {/* Peak */}
          <div className="flex flex-col items-center text-center space-y-1 sm:space-y-2 md:space-y-3 flex-shrink-0">
            <div className="w-12 h-12 sm:w-16 sm:h-16 md:w-20 md:h-20 lg:w-24 lg:h-24 rounded-lg sm:rounded-xl md:rounded-2xl bg-primary/20 flex items-center justify-center border-2 border-primary/40">
              <Mountain className="w-6 h-6 sm:w-8 sm:h-8 md:w-10 md:h-10 lg:w-12 lg:h-12 text-primary" />
            </div>
            <div>
              <div className="text-sm sm:text-base md:text-lg lg:text-xl xl:text-2xl font-bold text-primary">Peak</div>
              <div className="text-[10px] sm:text-xs md:text-sm text-muted-foreground hidden sm:block">from Peaker</div>
            </div>
          </div>

          {/* Equals */}
          <div className="flex items-center justify-center flex-shrink-0 -mb-0 sm:-mb-2 md:-mb-0">
            <div className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl xl:text-6xl font-bold text-muted-foreground">=</div>
          </div>

          {/* TerraPeak */}
          <div className="flex flex-col items-center text-center space-y-1 sm:space-y-2 md:space-y-3 flex-shrink-0">
            <div className="w-12 h-12 sm:w-16 sm:h-16 md:w-20 md:h-20 lg:w-24 lg:h-24 rounded-lg sm:rounded-xl md:rounded-2xl bg-gradient-to-br from-[#7B42BC]/20 to-primary/20 flex items-center justify-center border-2 border-primary">
              <Mountain className="w-6 h-6 sm:w-8 sm:h-8 md:w-10 md:h-10 lg:w-12 lg:h-12 text-primary" />
            </div>
            <div>
              <div className="text-sm sm:text-base md:text-lg lg:text-xl xl:text-2xl font-bold text-primary">TerraPeak</div>
              <div className="text-[10px] sm:text-xs md:text-sm text-muted-foreground hidden sm:block">Peak Performance</div>
            </div>
          </div>
        </div>

        <p className="text-center text-sm sm:text-base md:text-lg text-muted-foreground max-w-3xl mx-auto mt-6 sm:mt-8 md:mt-12 leading-relaxed px-2">
          Built by Peaker, TerraPeak combines the power of Terraform with peak-level caching performance to deliver the
          fastest registry proxy solution.
        </p>
      </div>
    </section>
  )
}
