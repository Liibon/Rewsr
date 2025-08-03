<div align="center">
  <img src="docs/images/logo.png" alt="Rewsr Logo" width="120" height="120">
  
  # Rewsr
  
  **CLI tool for AWS Nitro Enclaves**
  
  *Convert Docker images to EIF files and deploy them to AWS Nitro Enclaves with hardware attestation*
  
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Go Report Card](https://goreportcard.com/badge/github.com/rewsr/rewsr)](https://goreportcard.com/report/github.com/rewsr/rewsr)
  [![Release](https://img.shields.io/github/release/rewsr/rewsr.svg)](https://github.com/rewsr/rewsr/releases)
  
</div>

---

## 🚀 Quick Start

```bash
# Install
curl -sSL https://raw.githubusercontent.com/rewsr/rewsr/main/install.sh | bash

# Pack Docker image into EIF
rewsr pack nginx:alpine

# Deploy to Nitro Enclave
rewsr deploy nginx-alpine.eif --port 8080

# Generate attestation
rewsr attest nginx-alpine.eif

# Verify attestation
rewsr verify nginx-alpine.cbor
```

## 📋 Requirements

- AWS Nitro-enabled EC2 instance (c5n, m5n, r5n, t3, etc.)
- Docker installed
- AWS Nitro Enclaves CLI installed

## 📦 Installation

### Quick Install
```bash
curl -sSL https://raw.githubusercontent.com/rewsr/rewsr/main/install.sh | bash
```

### From Source
```bash
git clone https://github.com/rewsr/rewsr.git
cd rewsr
make install
```

### Download Binary
Download from [GitHub Releases](https://github.com/rewsr/rewsr/releases)

## 📚 Commands

<details>
<summary><strong>pack</strong> - Build EIF from Docker image</summary>

```bash
rewsr pack nginx:alpine
rewsr pack --output my-app.eif my-registry/app:latest
rewsr pack --entry '["python", "app.py"]' python:3.9
```

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output EIF filename | `{image}.eif` |
| `-e, --entry` | Override entrypoint | - |

</details>

<details>
<summary><strong>deploy</strong> - Launch in Nitro Enclave</summary>

```bash
rewsr deploy app.eif
rewsr deploy --port 8080 --cpu-count 4 --memory 4096 app.eif
```

| Flag | Description | Default |
|------|-------------|---------|
| `-p, --port` | Local port for vsock proxy | 0 |
| `--cpu-count` | Number of CPUs (1-16) | 2 |
| `--memory` | Memory in MB (512-16384) | 2048 |

</details>

<details>
<summary><strong>attest</strong> - Generate attestation document</summary>

```bash
rewsr attest app.eif
rewsr attest --output custom.cbor app.eif
```

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output CBOR file | `{name}.cbor` |

</details>

<details>
<summary><strong>verify</strong> - Verify attestation</summary>

```bash
rewsr verify app.cbor
```

</details>

## 🏢 Enterprise

For enterprise features including SSO, audit logging, compliance tooling, and 24/7 support:

**Contact:** [enterprise@rewsr.com](mailto:enterprise@rewsr.com)

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

[MIT License](LICENSE) for open source CLI.

Enterprise features require a separate commercial license.

## 🔗 Links

- **Issues:** [GitHub Issues](https://github.com/rewsr/rewsr/issues)
- **Enterprise:** [enterprise@rewsr.com](mailto:enterprise@rewsr.com)
- **Documentation:** Coming soon

---

<div align="center">
  <sub>Built with ❤️ for confidential computing</sub>
</div>
