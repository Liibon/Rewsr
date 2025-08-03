#!/bin/bash
# REWSR CLI installer script

set -e

# Configuration
REPO="rewsr/rewsr"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="rewsr"

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    "Linux") OS="linux" ;;
    "Darwin") OS="darwin" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

case $ARCH in
    "x86_64") ARCH="amd64" ;;
    "arm64"|"aarch64") ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Construct download URL
BINARY_NAME_FULL="${BINARY_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME_FULL="${BINARY_NAME_FULL}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME_FULL}"

echo "🔧 Installing REWSR CLI..."
echo "   OS: $OS"
echo "   Architecture: $ARCH"
echo "   Download URL: $DOWNLOAD_URL"

# Create install directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo "📁 Creating install directory: $INSTALL_DIR"
    sudo mkdir -p "$INSTALL_DIR"
fi

# Download and install
echo "⬇️  Downloading REWSR CLI..."
if command -v curl >/dev/null 2>&1; then
    sudo curl -sSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY_NAME}"
elif command -v wget >/dev/null 2>&1; then
    sudo wget -q "$DOWNLOAD_URL" -O "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# Verify installation
if [ -x "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    echo "✅ REWSR CLI installed successfully!"
    echo "🚀 Run 'rewsr --help' to get started."
    echo ""
    echo "Quick start:"
    echo "  rewsr pack nginx:alpine"
    echo "  rewsr deploy nginx:alpine-secure --port 8080"
    echo "  rewsr attest nginx:alpine-secure"
else
    echo "❌ Installation failed. Please check permissions and try again."
    exit 1
fi