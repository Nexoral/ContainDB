#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"
VERSION_FILE="./VERSION"
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')

# ensure only stable versions are published
if [[ "$VERSION" != *-stable ]]; then
  echo "❌ Stable version required to publish release. Current version: $VERSION"
  exit 0
fi

# collect all debs for this version
DEB_FILES=(./Packages/${APP_NAME}_${VERSION}_*) # to collect all files

TAG="v$VERSION"
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

# === Environment Variables ===
REPO="${GIT_REPOSITORY}" # GitHub Actions sets this automatically
TOKEN="${GIT_TOKEN}"     # GitHub Actions provides this

# === Build steps ===
./Scripts/BinBuilder.sh
echo "🔨 Binary Building completed of $APP_NAME version $VERSION for $ARCH"

./Scripts/PackageBuilder.sh
echo "📦 Package Building completed of $APP_NAME version $VERSION for $ARCH"

# === Checks ===
if [ ! -f "$VERSION_FILE" ]; then
  echo "❌ VERSION file not found"
  exit 1
fi

# ensure we have at least one .deb
if [ ${#DEB_FILES[@]} -eq 0 ]; then
  echo "❌ No .deb files found for version $VERSION in Packages/"
  exit 1
fi

if ! command -v gh &>/dev/null; then
  echo "❌ GitHub CLI (gh) not installed"

  # Update package list and install dependencies
  apt update
  apt install -y curl gnupg software-properties-common

  # Add GitHub CLI's official package repository
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg |
    dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg

  sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" |
    tee /etc/apt/sources.list.d/github-cli.list >/dev/null

  # Install gh
  apt update
  apt install -y gh

fi

# === Authenticate GitHub CLI ===
echo "🔑 Authenticating GitHub CLI..."
echo "${TOKEN}" | gh auth login --with-token

# === Create Release ===
echo "📦 Creating GitHub release for tag $TAG..."

gh release create "$TAG" "${DEB_FILES[@]}" \
  --title "$TAG" \
  --notes "🔨 Commit: $COMMIT_HASH

📝 Message:
$COMMIT_MSG"

echo "✅ GitHub release published with .deb assets"
