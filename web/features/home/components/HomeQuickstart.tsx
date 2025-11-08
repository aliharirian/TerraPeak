import React from "react";

const HomeQuickstart = () => {
  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch (err) {
      console.error("Failed to copy text:", err);
    }
  };
  return (
    <section id="quickstart" className="py-20 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
            Get Started in Minutes
          </h2>
          <p className="text-lg text-gray-400">
            Choose your preferred deployment method
          </p>
        </div>

        {/* Docker Compose */}
        <div className="mb-12">
          <h3 className="text-2xl font-semibold text-white mb-4 flex items-center gap-3">
            <span className="w-8 h-8 bg-green-500 text-black rounded-full flex items-center justify-center text-sm font-bold">
              1
            </span>
            Docker Compose (Recommended)
          </h3>
          <p className="text-gray-400 mb-4">
            The easiest way to get TerraPeak running with all dependencies:
          </p>

          <div className="code-block bg-gray-900 border border-gray-800 rounded-lg p-4 mb-4">
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs text-gray-500 font-mono">bash</span>
              <button
                className="copy-btn px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                onClick={() =>
                  copyToClipboard(
                    `git clone https://github.com/aliharirian/TerraPeak.git\ncd TerraPeak\ndocker-compose up -d`
                  )
                }
              >
                Copy
              </button>
            </div>
            <pre className="text-sm text-gray-300 overflow-x-auto">
              <code>{`git clone https://github.com/aliharirian/TerraPeak.git
cd TerraPeak
docker-compose up -d`}</code>
            </pre>
          </div>

          <p className="text-gray-400 text-sm">
            TerraPeak will be available on the configured ports. Check your{" "}
            <code className="px-2 py-0.5 bg-gray-800 text-green-400 rounded text-xs">
              cfg.yml
            </code>{" "}
            for server address settings.
          </p>
        </div>

        {/* Docker Run */}
        <div className="mb-12">
          <h3 className="text-2xl font-semibold text-white mb-4 flex items-center gap-3">
            <span className="w-8 h-8 bg-green-500 text-black rounded-full flex items-center justify-center text-sm font-bold">
              2
            </span>
            Docker Run
          </h3>
          <p className="text-gray-400 mb-4">
            Pull and run the latest TerraPeak image:
          </p>

          <div className="code-block bg-gray-900 border border-gray-800 rounded-lg p-4 mb-4">
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs text-gray-500 font-mono">bash</span>
              <button
                className="copy-btn px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                onClick={() =>
                  copyToClipboard("docker pull aliharirian/terrapeak:latest")
                }
              >
                Copy
              </button>
            </div>
            <pre className="text-sm text-gray-300 overflow-x-auto">
              <code>docker pull aliharirian/terrapeak:latest</code>
            </pre>
          </div>

          <div className="code-block bg-gray-900 border border-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs text-gray-500 font-mono">bash</span>
              <button
                className="copy-btn px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                onClick={() =>
                  copyToClipboard(
                    "docker run -d -p 8080:8080 -v $(pwd)/cfg.yml:/app/cfg.yml aliharirian/terrapeak:latest"
                  )
                }
              >
                Copy
              </button>
            </div>
            <pre className="text-sm text-gray-300 overflow-x-auto">
              <code>
                docker run -d -p 8080:8080 -v $(pwd)/cfg.yml:/app/cfg.yml
                aliharirian/terrapeak:latest
              </code>
            </pre>
          </div>
        </div>

        {/* Build from Source */}
        <div>
          <h3 className="text-2xl font-semibold text-white mb-4 flex items-center gap-3">
            <span className="w-8 h-8 bg-green-500 text-black rounded-full flex items-center justify-center text-sm font-bold">
              3
            </span>
            Build from Source
          </h3>
          <p className="text-gray-400 mb-4">
            For developers who want to build and customize:
          </p>

          <div className="code-block bg-gray-900 border border-gray-800 rounded-lg p-4">
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs text-gray-500 font-mono">bash</span>
              <button
                className="copy-btn px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded focus:outline-none focus:ring-2 focus:ring-green-500"
                onClick={() =>
                  copyToClipboard(
                    "cd registry\ngo build -o terrapeak\n./terrapeak"
                  )
                }
              >
                Copy
              </button>
            </div>
            <pre className="text-sm text-gray-300 overflow-x-auto">
              <code>{`cd registry
go build -o terrapeak
./terrapeak`}</code>
            </pre>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HomeQuickstart;
