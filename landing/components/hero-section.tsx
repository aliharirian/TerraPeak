import { Button } from "@/components/ui/button"
import { Github, Mountain, Play } from "lucide-react"
import Link from "next/link"

export function HeroSection() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden px-3 sm:px-4 md:px-6">
      {/* Background gradient effect */}
      <div className="absolute inset-0 bg-gradient-to-b from-primary/5 via-background to-background" />

      {/* Mountain peak decoration - hidden on mobile for better performance */}
      <div className="hidden sm:block absolute top-20 right-1/4 w-64 h-64 bg-primary/10 rounded-full blur-3xl animate-pulse-slow" />
      <div className="hidden sm:block absolute bottom-20 left-1/4 w-96 h-96 bg-primary/5 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: '2s' }} />

      <div className="relative max-w-6xl mx-auto text-center space-y-4 sm:space-y-6 md:space-y-8 py-12 sm:py-16 md:py-20">
        {/* Logo/Brand mark - Enhanced with animations */}
        <div className="flex items-center justify-center mb-4 sm:mb-6 md:mb-8">
          <div className="relative">
            {/* Animated glow effect behind logo - hidden on mobile */}
            <div className="hidden sm:block absolute inset-0 bg-primary/30 rounded-full blur-2xl animate-pulse-slow scale-150" />

            {/* Main logo with multiple animations */}
            <Mountain className="w-20 h-20 sm:w-28 sm:h-28 md:w-32 md:h-32 lg:w-40 lg:h-40 text-primary relative z-10 animate-float drop-shadow-[0_0_25px_rgba(134,239,172,0.5)]" />

            {/* Rotating ring effect - hidden on mobile */}
            <div className="hidden sm:block absolute inset-0 border-2 border-primary/20 rounded-full animate-spin-slow scale-125" />
            <div className="hidden sm:block absolute inset-0 border-2 border-primary/10 rounded-full animate-spin-reverse scale-150" />
          </div>
        </div>

        {/* Main headline */}
        <div className="space-y-2 sm:space-y-3 md:space-y-4">
          <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl xl:text-8xl font-bold tracking-tight text-balance animate-fade-in-up px-2">TerraPeak</h1>
          <p className="text-lg sm:text-xl md:text-2xl lg:text-3xl text-muted-foreground font-medium text-balance animate-fade-in-up animation-delay-200 px-2">
            Where Terraform Meets Peak Performance
          </p>
        </div>

        {/* Tagline */}
        <p className="text-base sm:text-lg md:text-xl lg:text-2xl text-foreground/80 max-w-3xl mx-auto leading-relaxed text-balance px-3">
          Cache Terraform providers locally. Download once, use forever.
          <span className="text-primary font-semibold"> 95% faster</span> builds with zero configuration.
        </p>

        {/* CTA Buttons */}
        <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-center gap-2.5 sm:gap-3 md:gap-4 pt-2 sm:pt-3 md:pt-4 px-3">
          <Button
            size="lg"
            className="text-base sm:text-lg px-6 py-5 sm:px-7 sm:py-5 md:px-8 md:py-6 bg-primary hover:bg-primary/90 text-primary-foreground w-full sm:w-auto"
            asChild
          >
            <Link href="#get-started-with-terrapeak">Get Started</Link>
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
          <Button
            size="lg"
            variant="ghost"
            className="text-base sm:text-lg px-6 py-5 sm:px-7 sm:py-5 md:px-8 md:py-6 hover:bg-primary/10 w-full sm:w-auto"
            asChild
          >
            <Link href="https://github.com/aliharirian/TerraPeak">
              <Github className="w-4 h-4 sm:w-5 sm:h-5 mr-2" />
              GitHub
            </Link>
          </Button>
        </div>

        {/* Quick stats banner */}
        <div className="grid grid-cols-3 gap-2 sm:gap-4 md:gap-6 max-w-4xl mx-auto pt-6 sm:pt-8 md:pt-12 px-3">
          <div className="bg-card/50 backdrop-blur border border-primary/20 rounded-lg p-2 sm:p-4 md:p-5 lg:p-6">
            <div className="text-xl sm:text-3xl md:text-4xl font-bold text-primary mb-1 sm:mb-2">5ms</div>
            <div className="text-[10px] sm:text-sm md:text-base text-muted-foreground">Response Time</div>
          </div>
          <div className="bg-card/50 backdrop-blur border border-primary/20 rounded-lg p-2 sm:p-4 md:p-5 lg:p-6">
            <div className="text-xl sm:text-3xl md:text-4xl font-bold text-primary mb-1 sm:mb-2">95%</div>
            <div className="text-[10px] sm:text-sm md:text-base text-muted-foreground">Bandwidth Savings</div>
          </div>
          <div className="bg-card/50 backdrop-blur border border-primary/20 rounded-lg p-2 sm:p-4 md:p-5 lg:p-6">
            <div className="text-xl sm:text-3xl md:text-4xl font-bold text-primary mb-1 sm:mb-2">{"<1s"}</div>
            <div className="text-[10px] sm:text-sm md:text-base text-muted-foreground">Sub-second Downloads</div>
          </div>
        </div>
      </div>
    </section>
  )
}
