"use client"

import type React from "react"
import {useState} from "react"

import {Button} from "@/components/ui/button"
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/components/ui/tabs"
import {AlertTriangle, Check, ChevronDown, ChevronUp, Copy} from "lucide-react"

function CodeBlock({code, language = "bash"}: { code: string; language?: string }) {
    const [copied, setCopied] = useState(false)

    const handleCopy = () => {
        navigator.clipboard.writeText(code)
        setCopied(true)
        setTimeout(() => setCopied(false), 2000)
    }

    return (
        <div className="bg-[#0D1117] border border-border rounded-md sm:rounded-lg overflow-hidden">
            <div className="flex items-center justify-between px-2.5 sm:px-3 md:px-4 py-1.5 sm:py-2 border-b border-border/50">
                <span className="text-[10px] sm:text-xs text-muted-foreground font-mono">{language}</span>
                <Button
                    size="sm"
                    variant="ghost"
                    className="hover:bg-primary/10 hover:text-primary h-6 sm:h-7 px-2 sm:px-3 transition-colors"
                    onClick={handleCopy}
                >
                    {copied ? (
                        <>
                            <Check className="w-2.5 h-2.5 sm:w-3 sm:h-3 text-primary mr-1"/>
                            <span className="text-[10px] sm:text-xs text-primary">Copied!</span>
                        </>
                    ) : (
                        <>
                            <Copy className="w-2.5 h-2.5 sm:w-3 sm:h-3 mr-1"/>
                            <span className="text-[10px] sm:text-xs">Copy</span>
                        </>
                    )}
                </Button>
            </div>
            <pre className="p-2.5 sm:p-3 md:p-4 overflow-x-auto text-xs sm:text-sm">
        <code className="font-mono text-green-400 leading-relaxed">{code}</code>
      </pre>
        </div>
    )
}

function CollapsibleSection({
                                title,
                                children,
                                defaultOpen = false,
                            }: {
    title: string
    children: React.ReactNode
    defaultOpen?: boolean
}) {
    const [isOpen, setIsOpen] = useState(defaultOpen)

    return (
        <div className="border border-border rounded-md sm:rounded-lg overflow-hidden">
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="w-full p-3 sm:p-4 bg-card hover:bg-card/80 transition-colors flex items-center justify-between"
            >
                <span className="font-bold text-left text-sm sm:text-base">{title}</span>
                {isOpen ? (
                    <ChevronUp className="w-4 h-4 sm:w-5 sm:h-5 text-primary flex-shrink-0"/>
                ) : (
                    <ChevronDown className="w-4 h-4 sm:w-5 sm:h-5 text-primary flex-shrink-0"/>
                )}
            </button>
            {isOpen && <div className="p-4 sm:p-5 md:p-6 bg-card/50 space-y-4 sm:space-y-5 md:space-y-6">{children}</div>}
        </div>
    )
}

export function QuickStart() {
    return (
        <section className="py-12 sm:py-16 md:py-20 lg:py-24 px-3 sm:px-4 md:px-6 bg-background" id="get-started-with-terrapeak">
            <div className="max-w-5xl mx-auto">
                <div className="text-center mb-6 sm:mb-8 md:mb-12">
                    <h2 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold mb-2 sm:mb-3 md:mb-4 text-balance px-2">Get Started with TerraPeak</h2>
                    <p className="text-base sm:text-lg text-muted-foreground max-w-2xl mx-auto px-2">Get TerraPeak running in under 2
                        minutes</p>
                </div>

                <div
                    className="mb-4 sm:mb-6 md:mb-8 p-3 sm:p-4 bg-yellow-500/10 border border-yellow-500/30 rounded-lg flex items-start gap-2 sm:gap-3">
                    <AlertTriangle className="w-4 h-4 sm:w-5 sm:h-5 text-yellow-500 flex-shrink-0 mt-0.5"/>
                    <div>
                        <h3 className="font-bold text-yellow-200 mb-0.5 sm:mb-1 text-sm sm:text-base">HTTPS Required for Production</h3>
                        <p className="text-xs sm:text-sm text-yellow-200/80 leading-relaxed">
                            Terraform requires HTTPS with valid SSL certificates to download providers. TerraPeak MUST
                            run behind an
                            HTTPS proxy (like Nginx with Let's Encrypt) for production use.{" "}
                            <code className="text-yellow-300 text-[10px] sm:text-xs">http://localhost:8081</code> only works for local testing.
                        </p>
                    </div>
                </div>

                <Tabs defaultValue="docker-compose" className="w-full">
                    <TabsList className="grid w-full grid-cols-2 h-auto gap-1 sm:gap-2 bg-transparent mb-4 sm:mb-6 md:mb-8">
                        <TabsTrigger
                            value="docker-compose"
                            className="data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:border-b-2 data-[state=active]:border-primary border-b-2 border-transparent flex items-center justify-center gap-1.5 sm:gap-2 md:gap-3 py-2.5 sm:py-3 md:py-4 transition-all"
                        >
                            <div className="text-center">
                                <div className="font-bold text-xs sm:text-sm md:text-base">Docker Compose</div>
                                <div className="text-[10px] sm:text-xs text-muted-foreground">Recommended</div>
                            </div>
                        </TabsTrigger>
                        <TabsTrigger
                            value="docker"
                            className="data-[state=active]:bg-primary/10 data-[state=active]:text-primary data-[state=active]:border-b-2 data-[state=active]:border-primary border-b-2 border-transparent flex items-center justify-center gap-1.5 sm:gap-2 md:gap-3 py-2.5 sm:py-3 md:py-4 transition-all"
                        >
                            <div className="text-center">
                                <div className="font-bold text-xs sm:text-sm md:text-base">Docker Command</div>
                                <div className="text-[10px] sm:text-xs text-muted-foreground">Quick & Simple</div>
                            </div>
                        </TabsTrigger>
                    </TabsList>

                    <TabsContent value="docker-compose" className="space-y-4 sm:space-y-5 md:space-y-6">
                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Clone and start</h3>
                            <CodeBlock
                                code={`git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak
docker-compose up -d

# Verify
curl http://localhost:8081/healthz`}
                            />
                            <p className="text-xs sm:text-sm text-muted-foreground mt-2 sm:mt-3">This starts TerraPeak on port 8081.</p>
                        </div>

                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Configure Terraform</h3>
                            <p className="text-xs sm:text-sm text-muted-foreground mb-2 sm:mb-3">
                                Update your Terraform configuration to use TerraPeak. Replace{" "}
                                <><strong className="text-primary">terraform.peaker.info</strong> with your actual
                                    domain
                                </>
                                :
                            </p>
                            <CodeBlock
                                language="hcl"
                                code={`# main.tf or provider.tf
terraform {
  required_providers {
    aws = {
      source  = "terraform.peaker.info/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`}
                            />
                        </div>

                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Run Terraform</h3>
                            <CodeBlock
                                code={`terraform init
terraform plan
terraform apply`}
                            />
                        </div>
                    </TabsContent>

                    <TabsContent value="docker" className="space-y-4 sm:space-y-5 md:space-y-6">
                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Download default config</h3>
                            <CodeBlock
                                code={`curl -O https://raw.githubusercontent.com/aliharirian/TerraPeak/main/registry/cfg.default.yml`}
                            />
                        </div>

                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Run container</h3>
                            <CodeBlock
                                code={`docker run -d \\
  --name terrapeak \\
  -p 8081:8081 \\
  -v $(pwd)/cfg.default.yml:/app/cfg.yml:ro \\
  aliharirian/terrapeak-registry:latest

# Verify
curl http://localhost:8081/healthz`}
                            />
                        </div>

                        <div>
                            <h3 className="font-bold mb-2 sm:mb-3 text-sm sm:text-base">Configure Terraform</h3>
                            <p className="text-xs sm:text-sm text-muted-foreground mb-2 sm:mb-3">
                                Update your Terraform configuration to use TerraPeak. Replace{" "}
                                <><strong className="text-primary">terraform.peaker.info</strong> with your actual
                                    domain
                                </>
                                :
                            </p>
                            <CodeBlock
                                language="hcl"
                                code={` # main.tf or provider.tf
terraform {
  required_providers {
    aws = {
      source  = "terraform.peaker.info/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`}
                            />
                        </div>
                    </TabsContent>
                </Tabs>

                <div className="mt-8 space-y-4">
                    <CollapsibleSection title="SSL/TLS Setup (Required for Production)">
                        <div>
                            <p className="text-sm text-muted-foreground mb-4">
                                Terraform requires HTTPS with valid SSL certificates. Setup with Nginx:
                            </p>
                            <CodeBlock
                                language="nginx"
                                code={`# /etc/nginx/sites-available/terrapeak
upstream terrapeak {
    server localhost:8081;
    keepalive 32;
}

server {
    listen 80;
    server_name terraform.peaker.info;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name terraform.peaker.info;

    # TLS Configuration
    ssl_certificate /etc/letsencrypt/live/terraform.peaker.info/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/terraform.peaker.info/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;

    # Proxy Configuration
    location / {
        proxy_pass http://terrapeak;
        proxy_http_version 1.1;
        
        # Preserve original request info
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Connection reuse
        proxy_set_header Connection "";
        
        # Timeouts for large provider downloads
        proxy_connect_timeout 60s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
        
        # Buffer settings
        proxy_buffering off;
        proxy_request_buffering off;
    }

    # Access logging
    access_log /var/log/nginx/terrapeak-access.log combined;
    error_log /var/log/nginx/terrapeak-error.log warn;
}`}
                            />
                        </div>

                        <div>
                            <h4 className="font-bold mb-3">Enable the configuration</h4>
                            <CodeBlock
                                code={`# Enable site
sudo ln -s /etc/nginx/sites-available/terrapeak /etc/nginx/sites-enabled/

# Get SSL certificate
sudo certbot --nginx -d terraform.peaker.info

# Test and reload Nginx
sudo nginx -t
sudo systemctl reload nginx`}
                            />
                        </div>
                    </CollapsibleSection>

                    <CollapsibleSection title="Troubleshooting">
                        <div>
                            <h4 className="font-bold mb-3">Connection Issues</h4>
                            <CodeBlock
                                code={`# Check service status
docker-compose ps
docker-compose logs -f terrapeak

# Verify health
curl http://localhost:8081/healthz`}
                            />
                        </div>

                        <div>
                            <h4 className="font-bold mb-3">Configuration Issues</h4>
                            <CodeBlock
                                code={`# Check config file is mounted correctly
docker exec terrapeak ls -la /app/cfg.yml

# View container logs
docker logs -f terrapeak`}
                            />
                        </div>

                        <div>
                            <h4 className="font-bold mb-3">Port Already in Use</h4>
                            <CodeBlock
                                code={`# Find what's using port 8081
lsof -i :8081

# Or use different port
docker run -p 9090:8081 ...`}
                            />
                        </div>

                        <div>
                            <h4 className="font-bold mb-3">Terraform Not Using TerraPeak</h4>
                            <CodeBlock
                                code={`# Verify your Terraform config uses correct domain
# Should be: terraform.peaker.info/hashicorp/aws

# Test the provider endpoint
curl https://terraform.peaker.info/v1/providers/hashicorp/aws/versions`}
                            />
                        </div>
                    </CollapsibleSection>
                </div>

                <div className="mt-12 p-8 bg-primary/10 border border-primary/30 rounded-xl">
                    <h3 className="text-2xl font-bold mb-4 flex items-center gap-2">
                        TerraPeak is running
                    </h3>

                    <div className="space-y-3 text-sm">
                        <p className="text-muted-foreground">Next steps:</p>
                        <ul className="space-y-2 text-muted-foreground">
                            <li className="flex items-start gap-2">
                                <span className="text-primary">→</span>
                                <span>Setup HTTPS with Nginx (required for production)</span>
                            </li>
                            <li className="flex items-start gap-2">
                                <span className="text-primary">→</span>
                                                <span>Update your Terraform configs to use terraform.peaker.info</span>
                            </li>
                            <li className="flex items-start gap-2">
                                <span className="text-primary">→</span>
                                <span>Run terraform init and watch your downloads get cached</span>
                            </li>
                        </ul>

                        <div className="pt-4 mt-4 border-t border-primary/30">
                            <p className="font-bold mb-2">Performance Benefits:</p>
                            <ul className="space-y-1 text-muted-foreground">
                                <li className="flex items-start gap-2">
                                    <span className="text-primary">•</span>
                                    <span>First download: ~2-5s (cached for next time)</span>
                                </li>
                                <li className="flex items-start gap-2">
                                    <span className="text-primary">•</span>
                                    <span>Subsequent downloads: ~5ms (95% faster!)</span>
                                </li>
                                <li className="flex items-start gap-2">
                                    <span className="text-primary">•</span>
                                    <span>Bandwidth savings: Up to 95%</span>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    )
}
