import React from "react";

const HomeConfiguration = () => {
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch (err) {
      console.error("Failed to copy text:", err);
    }
  };
  return (
    <section className="py-20 px-4 sm:px-6 lg:px-8 bg-gray-900/50">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Configuration
          </h2>
          <p className="text-lg text-gray-400">
            Configure TerraPeak with a simple YAML file
          </p>
        </div>

        <div className="bg-black border border-gray-800 rounded-xl p-6">
          <details className="group">
            <summary className="flex items-center justify-between text-lg font-semibold text-white mb-4 focus:outline-none focus:ring-2 focus:ring-green-500 rounded p-2">
              <span>Example cfg.yml Configuration</span>
              <svg
                className="w-5 h-5 text-gray-400 transition-transform group-open:rotate-180"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </summary>

            <div className="code-block bg-gray-900 border border-gray-800 rounded-lg p-4">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs text-gray-500 font-mono">yaml</span>
                <button
                  className="copy-btn px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                  onClick={() =>
                    copyToClipboard(`server:
  addr: ':8080'
  domain: 'https://terrapeak.example.com'  # Must be HTTPS

storage:
  type: 'minio'  # or 'file'
  minio:
    endpoint: 'minio:9000'
    accessKey: 'minioadmin'
    secretKey: 'minioadmin'
    bucket: 'terraform-providers'
    useSSL: false

proxy:
  enabled: true
  type: 'http'  # http, socks5, or socks4
  url: 'http://proxy.example.com:8080'`)
                  }
                >
                  Copy
                </button>
              </div>
              <pre className="text-sm text-gray-300 overflow-x-auto">
                <code>{`server:
  addr: ':8080'
  domain: 'https://terrapeak.example.com'  # Must be HTTPS

storage:
  type: 'minio'  # or 'file'
  minio:
    endpoint: 'minio:9000'
    accessKey: 'minioadmin'
    secretKey: 'minioadmin'
    bucket: 'terraform-providers'
    useSSL: false

proxy:
  enabled: true
  type: 'http'  # http, socks5, or socks4
  url: 'http://proxy.example.com:8080'`}</code>
              </pre>
            </div>
          </details>

          <div className="mt-6 p-4 bg-green-500/10 border border-green-500/20 rounded-lg">
            <div className="flex gap-3">
              <svg
                className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <p className="text-sm text-green-400 font-semibold mb-1">
                  Important: HTTPS Required
                </p>
                <p className="text-sm text-gray-300">
                  The{" "}
                  <code className="px-1.5 py-0.5 bg-gray-800 text-green-400 rounded text-xs">
                    server.domain
                  </code>{" "}
                  must use HTTPS with valid SSL certificates. Terraform will not
                  accept provider downloads over insecure connections.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HomeConfiguration;
