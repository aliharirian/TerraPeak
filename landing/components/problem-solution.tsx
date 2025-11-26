import { X, Check } from "lucide-react"

export function ProblemSolution() {
  const problems = [
    "Slow provider downloads",
    "Repeated bandwidth usage",
    "Upstream registry downtime",
    "Network bottlenecks",
  ]

  const solutions = ["Lightning-fast caching", "Smart storage backends", "Offline capability", "Instant retrieval"]

  return (
    <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-background">
      <div className="max-w-6xl mx-auto">
        <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-center mb-8 sm:mb-12 md:mb-16 text-balance px-2">
          Transform Your Terraform Workflow
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 md:gap-8">
          {/* Problem */}
          <div className="bg-card border border-destructive/30 rounded-xl sm:rounded-2xl p-4 sm:p-6 md:p-8">
            <h3 className="text-lg sm:text-xl md:text-2xl font-bold mb-4 sm:mb-5 md:mb-6 text-destructive">Without TerraPeak</h3>
            <ul className="space-y-2.5 sm:space-y-3 md:space-y-4">
              {problems.map((problem, index) => (
                <li key={index} className="flex items-start gap-2 sm:gap-2.5 md:gap-3">
                  <X className="w-5 h-5 sm:w-5.5 sm:h-5.5 md:w-6 md:h-6 text-destructive shrink-0 mt-0.5" />
                  <span className="text-sm sm:text-base md:text-lg text-muted-foreground">{problem}</span>
                </li>
              ))}
            </ul>
          </div>

          {/* Solution */}
          <div className="bg-card border border-primary/30 rounded-xl sm:rounded-2xl p-4 sm:p-6 md:p-8">
            <h3 className="text-lg sm:text-xl md:text-2xl font-bold mb-4 sm:mb-5 md:mb-6 text-primary">With TerraPeak</h3>
            <ul className="space-y-2.5 sm:space-y-3 md:space-y-4">
              {solutions.map((solution, index) => (
                <li key={index} className="flex items-start gap-2 sm:gap-2.5 md:gap-3">
                  <Check className="w-5 h-5 sm:w-5.5 sm:h-5.5 md:w-6 md:h-6 text-primary shrink-0 mt-0.5" />
                  <span className="text-sm sm:text-base md:text-lg text-foreground">{solution}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  )
}
