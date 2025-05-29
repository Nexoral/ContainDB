#!/bin/bash

# Options
FILE_TYPE=("deb" "tar.gz")
AVAILABLE_OPTIONS=("amd64" "arm64" "i386")
FOLDER_NAME="Packages"

for TYPE in "${FILE_TYPE[@]}"; do
  for ARCH in "${AVAILABLE_OPTIONS[@]}"; do
    echo "Building $TYPE package for architecture: $ARCH"

    # === CONFIG ===
    APP_NAME="containdb"
    # ARCH is set by outer loop
    MAINTAINER="Ankan Saha <ankansahaofficial@gmail.com>"
    DESCRIPTION="A short summary of what your program does."
    BINARY_PATH="./bin/ContainDB" # Path to your Go binary
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

    # Build per package type
    case "$TYPE" in
    deb)
      # === Create folder structure for .deb ===
      PKG_DIR="${APP_NAME}_${VERSION}_${ARCH}"
      mkdir -p "$PKG_DIR/DEBIAN"
      mkdir -p "$PKG_DIR/usr/local/bin"

      # === Write control file ===
      cat <<EOF >"$PKG_DIR/DEBIAN/control"
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
      if [ ! -d "$FOLDER_NAME" ]; then
        echo "Creating $FOLDER_NAME directory..."
        mkdir $FOLDER_NAME
      else
        echo "$FOLDER_NAME directory already exists."
      fi
      mv "${PKG_DIR}.deb" "./$FOLDER_NAME/${APP_NAME}_${VERSION}_${ARCH}.deb"
      rm -rf "$PKG_DIR"
      ;;
    tar.gz)
      # === Build tar.gz ===
      TMPDIR="$(mktemp -d)"
      mkdir -p "$TMPDIR/usr/local/bin"
      cp "$BINARY_PATH" "$TMPDIR/usr/local/bin/$APP_NAME"
      tar czf "./$FOLDER_NAME/${APP_NAME}_${VERSION}_${ARCH}.tar.gz" -C "$TMPDIR" .
      rm -rf "$TMPDIR"
      ;;
    esac

    echo "✅ $TYPE package created: ${APP_NAME}_${VERSION}_${ARCH}.${TYPE}"
  done
done
