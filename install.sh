#!/bin/sh
set -e

# Very basic VeFier installation script
REPO="wodorfells/vefier-cli"
BIN_DIR="/usr/local/bin"

echo "Downloading VeFier..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Replace with actual release fetcher in production
URL="https://github.com/wodorfells/vefier-cli/releases/latest/download/vefier_${OS}_${ARCH}.tar.gz"

curl -sL $URL | tar -xz -C /tmp
sudo mv /tmp/vefier $BIN_DIR/vefier

echo "✅ VeFier CLI успешно установлен!"
echo "Запустите 'vefier' для начала."