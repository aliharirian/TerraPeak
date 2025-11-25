export function IntegrationSection() {
  const integrations = [
    { name: "Terraform", color: "#7B42BC" },
    { name: "AWS", color: "#FF9900" },
    { name: "GCP", color: "#f44242" },
    { name: "Kubernetes", color: "#326CE5" },
    { name: "Azure", color: "#0078D4" },
    { name: "MongoDB Atlas", color: "#4DB33D" },
    { name: "Ali Cloud", color: "#FF6701" },
    { name: "Docker", color: "#2496ED" },
  ]

  return (
    <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-background">
      <div className="max-w-6xl mx-auto text-center">
        <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold mb-4 sm:mb-6 md:mb-8 text-balance px-2">Works With Your Stack</h2>
        <p className="text-base sm:text-lg md:text-xl text-muted-foreground mb-6 sm:mb-8 md:mb-12 max-w-3xl mx-auto px-2">
          Seamlessly integrates with all major cloud providers and storage backends
        </p>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-4 md:gap-5 lg:gap-6">
          {integrations.map((integration, index) => (
            <div
              key={index}
              className="bg-card border border-border rounded-lg sm:rounded-xl p-4 sm:p-5 md:p-6 lg:p-8 hover:border-primary/50 hover:scale-105 transition-all duration-300"
            >
              <div className="text-base sm:text-lg md:text-xl lg:text-2xl font-bold" style={{ color: integration.color }}>
                {integration.name}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
