import { Button } from "@/components/ui/button"
import { Github, Play } from "lucide-react"
import Link from "next/link"

export function CTASection() {
  return (
    <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-gradient-to-b from-primary/10 to-background">
      <div className="max-w-4xl mx-auto text-center space-y-4 sm:space-y-6 md:space-y-8">
        <h2 className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-balance px-2">Ready to Reach the Peak?</h2>
        <p className="text-base sm:text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto leading-relaxed px-3">
          Join developers who are accelerating their Terraform workflows with TerraPeak
        </p>

        <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-center gap-2.5 sm:gap-3 md:gap-4 pt-2 sm:pt-3 md:pt-4 px-3">
          <Button
            size="lg"
            className="text-base sm:text-lg px-6 py-5 sm:px-7 sm:py-5 md:px-8 md:py-6 bg-primary hover:bg-primary/90 text-primary-foreground w-full sm:w-auto"
            asChild
          >
            <Link href="https://github.com/aliharirian/TerraPeak">
              <Github className="w-4 h-4 sm:w-5 sm:h-5 mr-2" />
              Get Started Now
            </Link>
          </Button>
          <Button
            size="lg"
            variant="outline"
            className="text-base sm:text-lg px-6 py-5 sm:px-7 sm:py-5 md:px-8 md:py-6 border-primary/30 hover:bg-primary/10 bg-transparent w-full sm:w-auto"
            asChild
          >
            <Link href="/demo">
              <Play className="w-4 h-4 sm:w-5 sm:h-5 mr-2" />
              Demo
            </Link>
          </Button>
        </div>
      </div>
    </section>
  )
}
