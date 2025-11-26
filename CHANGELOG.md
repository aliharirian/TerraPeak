## [1.4.0](https://github.com/aliharirian/TerraPeak/compare/v1.3.1...v1.4.0) (2025-11-26)

### 🚀 Features

* add Docker support with multi-stage build and health check endpoint ([807248e](https://github.com/aliharirian/TerraPeak/commit/807248e7d0f452e3e6204b286a8cfaeb442b46ea))
* init landing page ([0b16bfe](https://github.com/aliharirian/TerraPeak/commit/0b16bfe9a649d6da2ec579c85a8c78dfae2e6d48))
* **landing:** add landing page ([d0dd92a](https://github.com/aliharirian/TerraPeak/commit/d0dd92a0c591ab5409db33be0c6fa57aacbf791e))

### 📚 Documentation

* add demo section and logo to README for improved user guidance ([46d140f](https://github.com/aliharirian/TerraPeak/commit/46d140f366f3fee4676b426592ed7450d74725d3))

### 👷 CI/CD

* simplify release configuration by removing unused prerelease branches ([c1eb634](https://github.com/aliharirian/TerraPeak/commit/c1eb634e5310269d133e1a12d523d0a7586ab996))
* update CI configuration to include landing page Docker build ([5d12f26](https://github.com/aliharirian/TerraPeak/commit/5d12f269b9acb1efbf0701825ddd124e07522c47))

## [1.3.1](https://github.com/aliharirian/TerraPeak/compare/v1.3.0...v1.3.1) (2025-11-24)

### 🚀 Features

* init bug report template for improved issue tracking ([9f8aed4](https://github.com/aliharirian/TerraPeak/commit/9f8aed493ffc0fd0daf745f160137b883acdecb6))

### 🐛 Bug Fixes

* **ci:** consolidate workflow improvements and release config updates ([ea7d71e](https://github.com/aliharirian/TerraPeak/commit/ea7d71ea036ad04b2d8ef7b68db71d75308c5025))
* **registry:** add docs, CI and registry fixes ([2db1f97](https://github.com/aliharirian/TerraPeak/commit/2db1f97910c827d0cfe933c7a6ed14ec843c2803))
* **registry:** update S3/MinIO configuration in .cfg.default.yml and docker-compose.yml for improved clarity and functionality & rename Terraform registry configuration to improve clarity ([7e2f907](https://github.com/aliharirian/TerraPeak/commit/7e2f9071d87643bbe4e91d659e33417f1c237d68))
* **tests:** rename proxy configuration fields for consistency and clarity ([27bdeb9](https://github.com/aliharirian/TerraPeak/commit/27bdeb99fc2e2013d42e13e1a19d065f22393de6))

### 📚 Documentation

* **registry:** write fully documentation ([0e4ee56](https://github.com/aliharirian/TerraPeak/commit/0e4ee56264655691453ec60989b8d75dac8adfc9))

# [1.3.0](https://github.com/aliharirian/TerraPeak/compare/v1.2.2...v1.3.0) (2025-11-21)


### Features

* **registry:** implement comprehensive Prometheus metrics instrumentation ([#5](https://github.com/aliharirian/TerraPeak/issues/5)) ([c8b5a3b](https://github.com/aliharirian/TerraPeak/commit/c8b5a3bb151e2c77342da0d194552e778e86fe9a))

## [1.2.2](https://github.com/aliharirian/TerraPeak/compare/v1.2.1...v1.2.2) (2025-10-22)


### Bug Fixes

* **registry:** consolidate cache proxy ([#2](https://github.com/aliharirian/TerraPeak/issues/2)) ([211b908](https://github.com/aliharirian/TerraPeak/commit/211b908fcae2c477ea38f0fd5bf47216f513e1ad))

## [1.2.1](https://github.com/aliharirian/TerraPeak/compare/v1.2.0...v1.2.1) (2025-10-20)


### Bug Fixes

* **registry:** implement interface store ([e6cfc00](https://github.com/aliharirian/TerraPeak/commit/e6cfc0079ef1c2512852e118574a5549a810f6a0))

# [1.2.0](https://github.com/aliharirian/TerraPeak/compare/v1.1.0...v1.2.0) (2025-09-30)


### Bug Fixes

* clean up proxy test files by removing unused imports and renaming variables for clarity ([1d36a95](https://github.com/aliharirian/TerraPeak/commit/1d36a951075eb4e6820045ebc100ea77142c663b))
* update proxy configuration structure in tests for improved clarity and maintainability ([73c78dd](https://github.com/aliharirian/TerraPeak/commit/73c78dd40bbcb4e0be9d9b37ba237e860d4eecb1))


### Features

* add proxy configuration to registry default YAML for enhanced network settings ([af1aa5b](https://github.com/aliharirian/TerraPeak/commit/af1aa5b5832fd8e0e8995e7e43920d97dd8b6bb5))
* implement proxy support in the registry, including configuration, handler, and client integration ([ebaad49](https://github.com/aliharirian/TerraPeak/commit/ebaad49e18a61bc77f8dc5f6454bc0492078c244))

# [1.1.0](https://github.com/aliharirian/TerraPeak/compare/v1.0.1...v1.1.0) (2025-09-18)


### Features

* add Nginx configuration with SSL support and Docker Compose setup for TerraPeak [skip ci] ([6c4f256](https://github.com/aliharirian/TerraPeak/commit/6c4f256ee828f12149c0ade4078bdf397b9bd15d))

## [1.0.1](https://github.com/aliharirian/TerraPeak/compare/v1.0.0...v1.0.1) (2025-09-18)


### Bug Fixes

* resolve Docker build issues with Go module downloads and Alpine base image ([9994749](https://github.com/aliharirian/TerraPeak/commit/9994749536fd2a655d2453ab9f19713ef1b447d0))

# 1.0.0 (2025-09-18)


### Bug Fixes

* correct typo in README roadmap section heading ([38c0d59](https://github.com/aliharirian/TerraPeak/commit/38c0d59d234a79e586c0dbc25418ac485da701e8))
