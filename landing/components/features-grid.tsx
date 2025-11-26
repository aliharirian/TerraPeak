import { Zap, Database, Lock, Globe, Activity, Settings } from "lucide-react"

export function FeaturesGrid() {
  const features = [
    {
      icon: Zap,
      title: "Lightning Fast",
      description: "Sub-second provider downloads from cache with 5ms response time",
    },
    {
      icon: Database,
      title: "Smart Storage",
      description: "Flexible GCS, Azure, or filesystem backends for reliable caching",
    },
    {
      icon: Lock,
      title: "Secure",
      description: "HTTPS/TLS support with SSL certificate configuration",
    },
    {
      icon: Globe,
      title: "Production Ready",
      description: "Deploy with Docker or standalone binaries",
    },
    {
      icon: Activity,
      title: "Observable",
      description: "Built-in health checks, metrics, and detailed logging",
    },
    {
      icon: Settings,
      title: "Configurable",
      description: "Flexible YAML configuration with upstream proxy support",
    },
  ]

  return (
    <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-background">
      <div className="max-w-6xl mx-auto">
        <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-center mb-8 sm:mb-12 md:mb-16 text-balance px-2">
          Everything You Need for Peak Performance
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4 md:gap-5 lg:gap-6 justify-items-center sm:justify-items-stretch">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-card border border-border rounded-lg sm:rounded-xl p-5 sm:p-6 md:p-7 lg:p-8 hover:border-primary/50 hover:bg-card/80 transition-all duration-300 group hover:shadow-lg hover:shadow-primary/10 w-full max-w-sm sm:max-w-none"
            >
              <feature.icon className="w-10 h-10 sm:w-11 sm:h-11 md:w-12 md:h-12 text-primary mb-3 sm:mb-4 group-hover:scale-110 transition-transform mx-auto sm:mx-0" />
              <h3 className="text-lg sm:text-xl font-bold mb-2 sm:mb-3 text-center sm:text-left">{feature.title}</h3>
              <p className="text-sm sm:text-base text-muted-foreground leading-relaxed text-center sm:text-left">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
