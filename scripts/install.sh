#!/bin/bash
# Installation script for gps CLI
# Usage: curl -sL https://raw.githubusercontent.com/wesbragagt/gps/main/scripts/install.sh | bash

set -e

REPO="wesbragagt/gps"
BINARY="gps"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux)  OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Determine extension
if [ "$OS" = "windows" ]; then
    EXT=".zip"
    BINARY="${BINARY}.exe"
else
    EXT=".tar.gz"
fi

# Get latest version from GitHub API
echo "Fetching latest release..."
LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
VERSION=$(curl -sL "$LATEST_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to get latest version" >&2
    exit 1
fi

echo "Latest version: ${VERSION}"

# Download URL
ARCHIVE_NAME="${BINARY}-${OS}-${ARCH}${EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

echo "Downloading ${ARCHIVE_NAME}..."
curl -sL -o "/tmp/${ARCHIVE_NAME}" "${DOWNLOAD_URL}"

# Extract
echo "Extracting..."
cd /tmp
if [ "$OS" = "windows" ]; then
    unzip -o "/tmp/${ARCHIVE_NAME}"
else
    tar -xzf "/tmp/${ARCHIVE_NAME}"
fi

# Install
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
echo "Installing to ${INSTALL_DIR}..."

if [ -w "${INSTALL_DIR}" ]; then
    mv "${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
    echo "Requires sudo to install to ${INSTALL_DIR}"
    sudo mv "${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"

# Cleanup
rm -f "/tmp/${ARCHIVE_NAME}"

echo ""
echo "Successfully installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
echo ""
echo "Run '${BINARY} --help' to get started."
