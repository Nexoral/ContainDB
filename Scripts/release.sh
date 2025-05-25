#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"
VERSION_FILE="./VERSION"
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
DEB_FILE="./Debian/${APP_NAME}_${VERSION}_${ARCH}.deb"
GITHUB_API="https://api.github.com"
REPO="${GIT_REPOSITORY}"  # GitHub Actions sets this automatically
TAG="v$VERSION"
TOKEN="${GIT_TOKEN}"      # GitHub Actions provides this
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

# Run other build scripts if needed
./Scripts/BinBuilder.sh
./Scripts/DebBuilder.sh

# === Safety checks ===
if [ ! -f "$VERSION_FILE" ]; then
  echo "❌ VERSION file not found"
  exit 1
fi

if [ ! -f "$DEB_FILE" ]; then
  echo "❌ .deb file not found at $DEB_FILE"
  exit 1
fi

if [ -z "$TOKEN" ]; then
  echo "❌ GITHUB_TOKEN is not set"
  exit 1
fi

# === Create GitHub Release ===
echo "📦 Creating GitHub release for tag $TAG"

CREATE_RESPONSE=$(curl -s -X POST "$GITHUB_API/repos/$REPO/releases" \
  -H "Authorization: token $TOKEN" \
  -H "Content-Type: application/json" \
  -d @- <<EOF
{
  "tag_name": "$TAG",
  "target_commitish": "$COMMIT_HASH",
  "name": "$TAG",
  "body": "🔨 Commit: $COMMIT_HASH\n\n📝 Message:\n$COMMIT_MSG",
  "draft": false,
  "prerelease": false
}
EOF
)

# === Extract upload URL ===
UPLOAD_URL=$(echo "$CREATE_RESPONSE" | grep '"upload_url":' | cut -d '"' -f 4 | sed "s/{?name,label}//")

if [ -z "$UPLOAD_URL" ]; then
  echo "❌ Failed to get upload URL"
  echo "Response: $CREATE_RESPONSE"
  exit 1
fi

# === Upload the .deb file ===
echo "📤 Uploading $(basename "$DEB_FILE")..."

curl -s --data-binary @"$DEB_FILE" \
  -H "Authorization: token $TOKEN" \
  -H "Content-Type: application/vnd.debian.binary-package" \
  "$UPLOAD_URL?name=$(basename $DEB_FILE)"

echo "✅ GitHub release published with .deb asset"
