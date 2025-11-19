import React from "react";

const HomeHero = () => {
  return (
    <section id="home" className="pt-32 pb-20 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <div className="text-center max-w-4xl mx-auto fade-in">
          <div className="inline-block mb-4 px-4 py-1.5 bg-green-500/10 border border-green-500/20 rounded-full">
            <span className="text-sm text-green-400 font-medium">
              Open Source • High Performance
            </span>
          </div>

          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-white mb-6 leading-tight">
            A high-performance caching proxy for Terraform Registry
          </h1>

          <p className="text-lg sm:text-xl text-gray-400 mb-8 leading-relaxed max-w-3xl mx-auto">
            Accelerate your Terraform provider downloads and reduce bandwidth
            costs. TerraPeak supports intelligent storage backends (MinIO, local
            file storage) and offers flexible proxy modes for corporate
            environments with outbound and inbound proxy support.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <a
              href="#quickstart"
              className="w-full sm:w-auto px-8 py-3.5 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black"
            >
              Get Started — Docker Compose
            </a>
            <a
              href="https://github.com/aliharirian/TerraPeak"
              target="_blank"
              rel="noopener noreferrer"
              className="w-full sm:w-auto px-8 py-3.5 bg-gray-900 hover:bg-gray-800 text-white font-semibold rounded-lg border border-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-black"
            >
              View on GitHub
            </a>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HomeHero;
