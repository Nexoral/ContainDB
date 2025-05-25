#!/bin/bash

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"  # Change to arm64 or other if needed
MAINTAINER="Ankan Saha <ankansahaofficial@gmail.com>"
DESCRIPTION="A short summary of what your program does."
BINARY_PATH="./bin/ContainDB"  # Path to your Go binary
VERSION_FILE="./VERSION"


# === Get version from VERSION file ===
if [ -f "$VERSION_FILE" ]; then
  VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
else
  echo "❌ VERSION file not found in project root"
  exit 1
fi


# === Check binary ===
if [ ! -f "$BINARY_PATH" ]; then
  echo "❌ Binary not found at $BINARY_PATH"
  exit 1
fi

# === Create folder structure ===
PKG_DIR="${APP_NAME}_${VERSION}_${ARCH}"
mkdir -p "$PKG_DIR/DEBIAN"
mkdir -p "$PKG_DIR/usr/local/bin"

# === Write control file ===
cat <<EOF > "$PKG_DIR/DEBIAN/control"
Package: $APP_NAME
Version: $VERSION
Section: utils
Priority: optional
Architecture: $ARCH
Maintainer: $MAINTAINER
Description: $DESCRIPTION
EOF

# === Copy binary ===
cp "$BINARY_PATH" "$PKG_DIR/usr/local/bin/$APP_NAME"
chmod 755 "$PKG_DIR/usr/local/bin/$APP_NAME"

# === Build the .deb ===
dpkg-deb --build "$PKG_DIR"

# === Move and clean up ===
mkdir Debian
mv "${PKG_DIR}.deb" "./Debian/${APP_NAME}_${VERSION}_${ARCH}.deb"
rm -rf "$PKG_DIR"

echo "✅ .deb package created: ${APP_NAME}_${VERSION}_${ARCH}.deb"
