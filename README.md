# Rewsr

**CLI tool for AWS Nitro Enclaves**

Rewsr converts Docker images to EIF files and deploys them to AWS Nitro Enclaves with hardware attestation.

## Quick Start

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

## Requirements

- AWS Nitro-enabled EC2 instance
- Docker
- AWS Nitro Enclaves CLI

## Installation

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

## Commands

### pack
Build EIF from Docker image:
```bash
rewsr pack nginx:alpine
rewsr pack --output my-app.eif my-registry/app:latest
rewsr pack --entry '["python", "app.py"]' python:3.9
```

### deploy
Launch in Nitro Enclave:
```bash
rewsr deploy app.eif
rewsr deploy --port 8080 --cpu-count 4 --memory 4096 app.eif
```

### attest
Generate attestation document:
```bash
rewsr attest app.eif
rewsr attest --output custom.cbor app.eif
```

### verify
Verify attestation:
```bash
rewsr verify app.cbor
```

## Enterprise

For enterprise features including SSO, audit logging, compliance tooling, and support:

**Contact:** enterprise@rewsr.com

## License

MIT License for open source CLI.

Enterprise features require commercial license.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Support

- Issues: [GitHub Issues](https://github.com/rewsr/rewsr/issues)
- Enterprise: enterprise@rewsr.com
