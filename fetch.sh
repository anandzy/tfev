#!/bin/bash

set -e

VERSION="SerialRelease2"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="amd64"
if [[ "$OS" == "darwin" ]]; then
  if sysctl -n machdep.cpu.brand_string | grep -q "Apple"; then
    ARCH="arm64"
  fi
fi

REPO="anandzy/tfev"
BINARY="filtervars"
URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY}-$OS-$ARCH"

echo "Downloading $BINARY version $VERSION for $OS/$ARCH..."
echo "From: $URL"

# Require GITHUB_TOKEN environment variable
if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "Error: Please export your GITHUB_TOKEN (GitHub personal access token) before running this script."
  exit 1
fi

curl -sL -H "Authorization: token $GITHUB_TOKEN" "$URL" -o "$BINARY"
chmod +x "$BINARY"

if ! ./"$BINARY" --version >/dev/null 2>&1; then
  echo "Error: downloaded file is not a valid executable."
  echo "Content of downloaded file:"
  cat "$BINARY"
  exit 1
fi