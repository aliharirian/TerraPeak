"use client"

import { Button } from "@/components/ui/button"
import { Footer } from "@/components/footer"
import { Mountain, Terminal, Copy, Check, Rocket, Zap, Shield, Home, Activity } from "lucide-react"
import Link from "next/link"
import { useState, useEffect } from "react"

export default function DemoPage() {
  const [copiedStep, setCopiedStep] = useState<number | null>(null)
  const [serviceStatus, setServiceStatus] = useState<'checking' | 'up' | 'down'>('checking')

  const copyToClipboard = (text: string, step: number) => {
    navigator.clipboard.writeText(text)
    setCopiedStep(step)
    setTimeout(() => setCopiedStep(null), 2000)
  }

  // Check service status
  useEffect(() => {
    const checkService = async () => {
      try {
        // Get health check URL from environment variable
        const healthCheckUrl = process.env.NEXT_PUBLIC_HEALTH_CHECK_URL || 'https://tf.tesaco.sbs/healthz'

        console.log('Checking health at:', healthCheckUrl)

        const controller = new AbortController()
        const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout

        const response = await fetch(healthCheckUrl, {
          method: 'GET',
          signal: controller.signal,
        })

        clearTimeout(timeoutId)
        console.log('Health check response status:', response.status)

        // Check if response is ok (200-299)
        setServiceStatus(response.ok ? 'up' : 'down')
      } catch (error) {
        console.error('Health check failed:', error)
        setServiceStatus('down')
      }
    }

    checkService()
    const interval = setInterval(checkService, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Transparent Fixed Header */}
      <header className="fixed top-0 left-0 right-0 z-50 border-b border-border/50 bg-background/80 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 sm:py-4 flex items-center justify-between">
          {/* Left: Home Icon */}
          <Link
            href="/"
            className="flex items-center gap-2 sm:gap-3 text-foreground/80 hover:text-primary transition-colors group"
            aria-label="Go to homepage"
          >
            <div className="p-2 sm:p-2.5 rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-all duration-300">
              <Home className="w-4 h-4 sm:w-5 sm:h-5 text-primary" />
            </div>
            <span className="font-semibold text-xs sm:text-sm hidden xs:inline">Home</span>
          </Link>

          {/* Right: Service Status */}
          <div className={`flex items-center gap-1.5 sm:gap-2 px-2.5 sm:px-4 py-1.5 sm:py-2 rounded-lg backdrop-blur-sm ${
            serviceStatus === 'checking' ? 'bg-yellow-500/20 border border-yellow-500/50' :
            serviceStatus === 'up' ? 'bg-green-500/20 border border-green-500/50' :
            'bg-red-500/20 border border-red-500/50'
          }`}>
            <Activity className={`w-3.5 h-3.5 sm:w-4 sm:h-4 ${
              serviceStatus === 'checking' ? 'text-yellow-500' :
              serviceStatus === 'up' ? 'text-green-500' :
              'text-red-500'
            }`} />
            <span className={`text-xs sm:text-sm font-medium ${
              serviceStatus === 'checking' ? 'text-yellow-500' :
              serviceStatus === 'up' ? 'text-green-500' :
              'text-red-500'
            }`}>
              {serviceStatus === 'checking' ? 'Checking' :
               serviceStatus === 'up' ? 'Service Up' :
               'Service Down'}
            </span>
            <div className={`w-1.5 h-1.5 sm:w-2 sm:h-2 rounded-full ${
              serviceStatus === 'checking' ? 'bg-yellow-500 animate-pulse' :
              serviceStatus === 'up' ? 'bg-green-500 animate-pulse' :
              'bg-red-500'
            }`} />
          </div>
        </div>
      </header>

      {/* Content with padding for fixed header */}
      <div className="pt-16 sm:pt-20">
        {/* Hero Section */}
        <section className="relative flex items-center justify-center overflow-hidden px-3 sm:px-4 py-8 sm:py-12 md:py-16 lg:py-20 bg-gradient-to-b from-primary/5 via-background to-background">
        {/* Animated background orbs - hidden on mobile for better performance */}
        <div className="hidden sm:block absolute top-10 sm:top-20 right-1/4 w-32 h-32 sm:w-64 sm:h-64 bg-primary/10 rounded-full blur-3xl animate-pulse" />
        <div className="hidden sm:block absolute bottom-10 sm:bottom-20 left-1/4 w-48 h-48 sm:w-96 sm:h-96 bg-primary/5 rounded-full blur-3xl animate-pulse delay-1000" />

        <div className="relative max-w-5xl mx-auto text-center space-y-3 sm:space-y-4 md:space-y-6">
          {/* Logo */}
          <div className="flex items-center justify-center mb-1 sm:mb-2 md:mb-4 animate-fade-in">
            <Mountain className="w-10 h-10 sm:w-12 sm:h-12 md:w-16 md:h-16 text-primary animate-bounce-slow" />
          </div>

          {/* Headline */}
          <div className="space-y-1.5 sm:space-y-2 md:space-y-3 animate-slide-up">
            <h1 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl xl:text-6xl font-bold tracking-tight text-balance px-2">
              Try TerraPeak Now
            </h1>
            <p className="text-sm sm:text-base md:text-lg lg:text-xl text-muted-foreground font-medium text-balance px-2">
              Get started in less than 2 minutes
            </p>
          </div>

          {/* Features badges */}
          <div className="flex flex-wrap items-center justify-center gap-1.5 sm:gap-2 md:gap-3 pt-1 sm:pt-2 animate-fade-in-delay px-2">
            <div className="flex items-center gap-1 sm:gap-1.5 md:gap-2 bg-card/50 backdrop-blur border border-primary/20 rounded-full px-2 py-1 sm:px-3 sm:py-1.5 md:px-4 md:py-2">
              <Rocket className="w-3 h-3 sm:w-3.5 sm:h-3.5 md:w-4 md:h-4 text-primary" />
              <span className="text-[10px] sm:text-xs md:text-sm font-medium">Fast Setup</span>
            </div>
            <div className="flex items-center gap-1 sm:gap-1.5 md:gap-2 bg-card/50 backdrop-blur border border-primary/20 rounded-full px-2 py-1 sm:px-3 sm:py-1.5 md:px-4 md:py-2">
              <Zap className="w-3 h-3 sm:w-3.5 sm:h-3.5 md:w-4 md:h-4 text-primary" />
              <span className="text-[10px] sm:text-xs md:text-sm font-medium">Instant Cache</span>
            </div>
            <div className="flex items-center gap-1 sm:gap-1.5 md:gap-2 bg-card/50 backdrop-blur border border-primary/20 rounded-full px-2 py-1 sm:px-3 sm:py-1.5 md:px-4 md:py-2">
              <Shield className="w-3 h-3 sm:w-3.5 sm:h-3.5 md:w-4 md:h-4 text-primary" />
              <span className="text-[10px] sm:text-xs md:text-sm font-medium">Secure</span>
            </div>
          </div>
        </div>
      </section>

      {/* Instructions Section */}
      <section className="py-6 sm:py-8 md:py-12 lg:py-16 px-3 sm:px-4 bg-background">
        <div className="max-w-4xl mx-auto">
          <div className="space-y-3 sm:space-y-4 md:space-y-6">
            {/* Step 1 */}
            <div className="group bg-card border border-border rounded-lg sm:rounded-xl md:rounded-2xl p-3 sm:p-4 md:p-6 lg:p-8 hover:shadow-lg transition-all duration-300 hover:scale-[1.01] hover:border-primary/30">
              <div className="flex items-start gap-2 sm:gap-3 md:gap-4">
                {/* Step number - hidden on mobile, visible on sm+ */}
                <div className="hidden sm:flex flex-shrink-0 w-8 h-8 md:w-10 md:h-10 rounded-full bg-primary/10 items-center justify-center text-primary font-bold text-base md:text-lg">
                  1
                </div>
                <div className="flex-1 space-y-2 sm:space-y-3 md:space-y-4 min-w-0">
                  <h3 className="text-lg sm:text-xl md:text-2xl font-semibold">Update Your Terraform Configuration</h3>
                  <p className="text-xs sm:text-sm md:text-base text-muted-foreground leading-relaxed">
                    Modify your main Terraform file or provider configuration file to use TerraPeak as the provider source. Replace{" "}
                    <code className="px-1.5 py-0.5 sm:px-2 sm:py-1 bg-primary/10 rounded text-primary font-mono text-[10px] sm:text-xs md:text-sm">
                      terraform.peaker.info
                    </code>{" "}
                    with your actual TerraPeak domain.
                  </p>

                  <div className="relative group/code">
                    <div className="absolute top-1.5 right-1.5 sm:top-2 sm:right-2 md:top-3 md:right-3 z-10">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 p-0 opacity-0 group-hover/code:opacity-100 transition-opacity"
                        onClick={() =>
                          copyToClipboard(
                            `# main.tf or provider.tf
terraform {
  required_providers {
    aws = {
      source  = "terraform.peaker.info/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`,
                            1
                          )
                        }
                      >
                        {copiedStep === 1 ? (
                          <Check className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4 text-green-500" />
                        ) : (
                          <Copy className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4" />
                        )}
                      </Button>
                    </div>
                    <pre className="bg-muted/50 border border-border rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4 overflow-x-auto">
                      <code className="text-[10px] sm:text-xs md:text-sm font-mono">
                        <span className="text-muted-foreground"># main.tf or provider.tf</span>
                        {"\n"}
                        <span className="text-blue-500">terraform</span> {"{"}
                        {"\n  "}
                        <span className="text-blue-500">required_providers</span> {"{"}
                        {"\n    "}
                        <span className="text-blue-500">aws</span> = {"{"}
                        {"\n      "}
                        <span className="text-blue-500">source</span>  ={" "}
                        <span className="text-green-500">"terraform.peaker.info/hashicorp/aws"</span>
                        {"\n      "}
                        <span className="text-blue-500">version</span> ={" "}
                        <span className="text-green-500">"~{'>'}5.0"</span>
                        {"\n    "}
                        {"}"}
                        {"\n  "}
                        {"}"}
                        {"\n}"}
                      </code>
                    </pre>
                  </div>

                  <div className="bg-blue-200 dark:bg-blue-950/60 border border-blue-400 dark:border-blue-600 rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4">
                    <p className="text-[10px] sm:text-xs md:text-sm text-blue-900 dark:text-blue-100 leading-relaxed">
                      <strong>💡 Pro Tip:</strong> Replace{" "}
                      <code className="px-1 py-0.5 bg-blue-300 dark:bg-blue-900/80 rounded text-[10px] font-mono">
                        terraform.peaker.info
                      </code>{" "}
                      with your actual TerraPeak domain. You can configure multiple providers this way.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Step 2 */}
            <div className="group bg-card border border-border rounded-lg sm:rounded-xl md:rounded-2xl p-3 sm:p-4 md:p-6 lg:p-8 hover:shadow-lg transition-all duration-300 hover:scale-[1.01] hover:border-primary/30">
              <div className="flex items-start gap-2 sm:gap-3 md:gap-4">
                {/* Step number - hidden on mobile, visible on sm+ */}
                <div className="hidden sm:flex flex-shrink-0 w-8 h-8 md:w-10 md:h-10 rounded-full bg-primary/10 items-center justify-center text-primary font-bold text-base md:text-lg">
                  2
                </div>
                <div className="flex-1 space-y-2 sm:space-y-3 md:space-y-4 min-w-0">
                  <h3 className="text-lg sm:text-xl md:text-2xl font-semibold">Initialize Terraform</h3>
                  <p className="text-xs sm:text-sm md:text-base text-muted-foreground leading-relaxed">
                    Run <code className="px-1 sm:px-1.5 md:px-2 py-0.5 sm:py-1 bg-primary/10 rounded text-primary font-mono text-[10px] sm:text-xs md:text-sm">terraform init</code> to
                    initialize your Terraform workspace with TerraPeak.
                  </p>

                  <div className="relative group/code">
                    <div className="absolute top-1.5 right-1.5 sm:top-2 sm:right-2 md:top-3 md:right-3 z-10">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 p-0 opacity-0 group-hover/code:opacity-100 transition-opacity"
                        onClick={() => copyToClipboard("terraform init", 2)}
                      >
                        {copiedStep === 2 ? (
                          <Check className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4 text-green-500" />
                        ) : (
                          <Copy className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4" />
                        )}
                      </Button>
                    </div>
                    <pre className="bg-muted/50 border border-border rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4 overflow-x-auto flex items-center gap-1.5 sm:gap-2 md:gap-3">
                      <Terminal className="w-3.5 h-3.5 sm:w-4 sm:h-4 md:w-5 md:h-5 text-primary flex-shrink-0" />
                      <code className="text-[10px] sm:text-xs md:text-sm font-mono">
                        <span className="text-primary">terraform init</span>
                      </code>
                    </pre>
                  </div>
                </div>
              </div>
            </div>

            {/* Step 3 */}
            <div className="group bg-card border border-border rounded-lg sm:rounded-xl md:rounded-2xl p-3 sm:p-4 md:p-6 lg:p-8 hover:shadow-lg transition-all duration-300 hover:scale-[1.01] hover:border-primary/30">
              <div className="flex items-start gap-2 sm:gap-3 md:gap-4">
                {/* Step number - hidden on mobile, visible on sm+ */}
                <div className="hidden sm:flex flex-shrink-0 w-8 h-8 md:w-10 md:h-10 rounded-full bg-primary/10 items-center justify-center text-primary font-bold text-base md:text-lg">
                  3
                </div>
                <div className="flex-1 space-y-2 sm:space-y-3 md:space-y-4 min-w-0">
                  <h3 className="text-lg sm:text-xl md:text-2xl font-semibold">Plan & Apply</h3>
                  <p className="text-xs sm:text-sm md:text-base text-muted-foreground leading-relaxed">
                    Run your Terraform commands as usual. TerraPeak will automatically cache and accelerate your provider downloads.
                  </p>

                  <div className="space-y-2 sm:space-y-3">
                    <div className="relative group/code">
                      <div className="absolute top-1.5 right-1.5 sm:top-2 sm:right-2 md:top-3 md:right-3 z-10">
                        <Button
                          size="sm"
                          variant="ghost"
                          className="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 p-0 opacity-0 group-hover/code:opacity-100 transition-opacity"
                          onClick={() => copyToClipboard("terraform plan", 3)}
                        >
                          {copiedStep === 3 ? (
                            <Check className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4 text-green-500" />
                          ) : (
                            <Copy className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4" />
                          )}
                        </Button>
                      </div>
                      <pre className="bg-muted/50 border border-border rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4 overflow-x-auto flex items-center gap-1.5 sm:gap-2 md:gap-3">
                        <Terminal className="w-3.5 h-3.5 sm:w-4 sm:h-4 md:w-5 md:h-5 text-primary flex-shrink-0" />
                        <code className="text-[10px] sm:text-xs md:text-sm font-mono">
                          <span className="text-primary">terraform plan</span>
                        </code>
                      </pre>
                    </div>

                    <div className="relative group/code">
                      <div className="absolute top-1.5 right-1.5 sm:top-2 sm:right-2 md:top-3 md:right-3 z-10">
                        <Button
                          size="sm"
                          variant="ghost"
                          className="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 p-0 opacity-0 group-hover/code:opacity-100 transition-opacity"
                          onClick={() => copyToClipboard("terraform apply", 4)}
                        >
                          {copiedStep === 4 ? (
                            <Check className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4 text-green-500" />
                          ) : (
                            <Copy className="h-3 w-3 sm:h-3.5 sm:w-3.5 md:h-4 md:w-4" />
                          )}
                        </Button>
                      </div>
                      <pre className="bg-muted/50 border border-border rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4 overflow-x-auto flex items-center gap-1.5 sm:gap-2 md:gap-3">
                        <Terminal className="w-3.5 h-3.5 sm:w-4 sm:h-4 md:w-5 md:h-5 text-primary flex-shrink-0" />
                        <code className="text-[10px] sm:text-xs md:text-sm font-mono">
                          <span className="text-primary">terraform apply</span>
                        </code>
                      </pre>
                    </div>
                  </div>

                  <div className="bg-green-200 dark:bg-green-950/60 border border-green-400 dark:border-green-600 rounded-md sm:rounded-lg p-2 sm:p-3 md:p-4">
                    <p className="text-[10px] sm:text-xs md:text-sm text-green-900 dark:text-green-100 leading-relaxed">
                      <strong>🎉 That's it!</strong> Your Terraform downloads are now accelerated by TerraPeak. Enjoy faster builds
                      and reduced bandwidth usage.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* CTA Section */}
          <div className="mt-6 sm:mt-8 md:mt-12 text-center space-y-2 sm:space-y-3 md:space-y-4 bg-gradient-to-br from-primary/10 via-primary/5 to-transparent border border-primary/20 rounded-lg sm:rounded-xl md:rounded-2xl p-4 sm:p-6 md:p-8 lg:p-10">
            <h3 className="text-xl sm:text-2xl md:text-3xl font-bold">Ready to Accelerate?</h3>
            <p className="text-xs sm:text-sm md:text-base lg:text-lg text-muted-foreground max-w-2xl mx-auto px-2 leading-relaxed">
              Deploy your own TerraPeak instance or explore the full documentation to learn more.
            </p>
            <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-center gap-2 sm:gap-3 md:gap-4 pt-1 sm:pt-2">
              <Button size="lg" className="text-xs sm:text-sm md:text-base px-4 py-5 sm:px-5 sm:py-5 md:px-6 md:py-6 w-full sm:w-auto" asChild>
                <Link href="/">Get Started</Link>
              </Button>
              <Button size="lg" variant="outline" className="text-xs sm:text-sm md:text-base px-4 py-5 sm:px-5 sm:py-5 md:px-6 md:py-6 w-full sm:w-auto" asChild>
                <Link href="https://github.com/aliharirian/TerraPeak">View on GitHub</Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <Footer />

      <style jsx global>{`
        @keyframes fade-in {
          from {
            opacity: 0;
          }
          to {
            opacity: 1;
          }
        }

        @keyframes slide-up {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        @keyframes bounce-slow {
          0%,
          100% {
            transform: translateY(0);
          }
          50% {
            transform: translateY(-10px);
          }
        }

        .animate-fade-in {
          animation: fade-in 1s ease-out;
        }

        .animate-fade-in-delay {
          animation: fade-in 1s ease-out 0.3s both;
        }

        .animate-slide-up {
          animation: slide-up 1s ease-out 0.2s both;
        }

        .animate-bounce-slow {
          animation: bounce-slow 3s ease-in-out infinite;
        }

        .delay-1000 {
          animation-delay: 1s;
        }
      `}</style>
      </div>
    </div>
  )
}

