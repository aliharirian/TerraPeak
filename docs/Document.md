# TerraPeak Documentation

Welcome to the TerraPeak documentation! This page serves as a guide to help you navigate our comprehensive documentation.

## 📚 Documentation Structure

TerraPeak's documentation is organized into the following sections:

### Main Documentation

- **[README](../README.md)** - Quick start guide and overview
- **[INSTALLATION](INSTALLATION.md)** - Detailed installation instructions for all deployment methods
- **[CONTRIBUTING](../CONTRIBUTING.md)** - Guidelines for contributing to TerraPeak
- **[LICENSE](../LICENSE)** - Apache License 2.0

### Registry Documentation

Located in [`docs/registry/`](registry/):

- **[API Documentation](registry/API.md)** - Complete API reference and endpoint documentation
- **[Architecture](registry/ARCHITECTURE.md)** - System design, component details, and technical architecture

## 🚀 Quick Navigation

### Getting Started
1. Start with the [README](../README.md) for a quick overview
2. Follow the [INSTALLATION guide](INSTALLATION.md) to deploy TerraPeak
3. Configure your Terraform to use TerraPeak

### For Developers
1. Read the [ARCHITECTURE](registry/ARCHITECTURE.md) documentation
2. Review the [CONTRIBUTING](../CONTRIBUTING.md) guidelines
3. Check the [API Documentation](registry/API.md) for endpoint details

### For API Users
1. Review [API Documentation](registry/API.md) for all available endpoints
2. Check the configuration section in the [README](../README.md)
3. See usage examples and troubleshooting

## 📖 Documentation Topics

### Installation & Deployment
- [Docker Compose Setup](INSTALLATION.md#docker-compose-recommended)
- [Docker Installation](INSTALLATION.md#docker)
- [Binary Installation](INSTALLATION.md#binary-installation)
- [Kubernetes Deployment](INSTALLATION.md#kubernetes-deployment)
- [Building from Source](INSTALLATION.md#build-from-source)

### Configuration
- [Basic Configuration](../README.md#configuration)
- [Storage Backends](registry/ARCHITECTURE.md#storage-architecture)
- [Proxy Configuration](INSTALLATION.md#run-with-sproxy-backend)
- [SSL/TLS Setup](../README.md#ssltls-setup)

### API Reference
- [Registry API Endpoints](registry/API.md#registry-api)
- [Proxy API Endpoints](registry/API.md#proxy-api)
- [Error Handling](registry/API.md#error-handling)
- [Response Format](registry/API.md#response-format)

### Architecture & Design
- [System Overview](registry/ARCHITECTURE.md#system-architecture)
- [Component Design](registry/ARCHITECTURE.md#component-design)
- [Storage Architecture](registry/ARCHITECTURE.md#storage-architecture)
- [Caching Strategy](registry/ARCHITECTURE.md#caching-strategy)
- [Performance Optimization](registry/ARCHITECTURE.md#performance-optimization)

### Development
- [Development Setup](../CONTRIBUTING.md#getting-started)
- [Development Workflow](../CONTRIBUTING.md#development-workflow)
- [Testing Guidelines](../CONTRIBUTING.md#testing)
- [Coding Standards](../CONTRIBUTING.md#coding-standards)

## 🔍 Need Help?

- **Bug Reports**: [Create an issue](https://github.com/aliharirian/TerraPeak/issues/new)
- **Questions**: Use [GitHub Discussions](https://github.com/aliharirian/TerraPeak/discussions)
- **General Support**: See [README - Support section](../README.md#support)

## 🗺️ Document Map

```
TerraPeak/
├── README.md                    # Main project overview
├── CONTRIBUTING.md              # Contribution guidelines
├── LICENSE                      # Apache 2.0 License
├── CHANGELOG.md                 # Version history
├── docs/
│   ├── Document.md             # This navigation guide
│   ├── INSTALLATION.md         # Installation guide
│   └── registry/
│       ├── API.md              # API reference
│       └── ARCHITECTURE.md     # Architecture documentation
└── registry/                    # Source code
    └── [Go source files]
```

## 📝 Contributing to Documentation

Documentation improvements are always welcome! If you find any issues or want to improve the documentation:

1. Fork the repository
2. Make your changes
3. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

---

**Last Updated**: November 2025
**Version**: 1.0+

