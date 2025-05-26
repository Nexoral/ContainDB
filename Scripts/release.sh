#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"
VERSION_FILE="./VERSION"
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
DEB_FILE="./Debian/${APP_NAME}_${VERSION}_${ARCH}.deb"
TAG="v$VERSION"
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

# === Environment Variables ===
REPO="${GIT_REPOSITORY}"  # GitHub Actions sets this automatically
TOKEN="${GIT_TOKEN}"      # GitHub Actions provides this

# === Build steps ===
./Scripts/BinBuilder.sh
./Scripts/DebBuilder.sh

# === Checks ===
if [ ! -f "$VERSION_FILE" ]; then
  echo "❌ VERSION file not found"
  exit 1
fi

if [ ! -f "$DEB_FILE" ]; then
  echo "❌ .deb file not found at $DEB_FILE"
  exit 1
fi

if ! command -v gh &> /dev/null; then
  echo "❌ GitHub CLI (gh) not installed"

  # Update package list and install dependencies
  apt update
  apt install -y curl gnupg software-properties-common

  # Add GitHub CLI's official package repository
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | \
    dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg

  sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | \
    tee /etc/apt/sources.list.d/github-cli.list > /dev/null

  # Install gh
  apt update
  apt install -y gh


fi

# === Create Release ===
echo "📦 Creating GitHub release for tag $TAG..."


 echo -e "${CYAN}📡 Creating GitHub release...${RESET}"

    RELEASE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "Accept: application/vnd.github+json" \
        https://api.github.com/repos/AnkanSaha/${REPO}/releases \
        -d "{
            \"tag_name\": \"v$VERSION\",
            \"target_commitish\": \"main\",
            \"name\": \"v$VERSION\",
            \"body\": \"Auto-generated release for 🔨 Commit: $COMMIT_HASH\",
            \"draft\": false,
            \"prerelease\": false
        }")

    UPLOAD_URL=$(echo "$RELEASE_RESPONSE" | grep upload_url | cut -d '"' -f 4 | cut -d '{' -f 1)

    if [ -z "$UPLOAD_URL" ]; then
        echo -e "${RED}❌ Failed to create GitHub release."
        echo "$RELEASE_RESPONSE"
        exit 1
    fi

    echo -e "${CYAN}📦 Uploading .deb to GitHub release..."
    curl -s -X POST \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "Content-Type: application/vnd.debian.binary-package" \
        --data-binary @"${DEB_FILE} \
        ""${UPLOAD_URL}"?name=${APP_NAME}"

    echo -e "🎉 Uploaded ${APP_NAME} to GitHub Releases successfully!"

📝 Message:
$COMMIT_MSG"

echo "✅ GitHub release published with .deb asset"
