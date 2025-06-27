VERSION="feature-feature%2F1a-3"  # You can make this dynamic!
REPO="anandzy/tfev"
BINARY="filtervars"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Download URL
URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY}-${OS}-${ARCH}"

echo "Downloading $BINARY version $VERSION for $OS/$ARCH..."
curl -L "$URL" -o "$BINARY"
chmod +x "$BINARY"
sudo mv "$BINARY" /usr/local/bin/


curl -L -v "https://github.com/anandzy/tfev/releases/download/$VERSION/filtervars-darwin-amd64" -o filtervars
file filtervars
head filtervars