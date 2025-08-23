#!/bin/bash
set -e


# Check which package is needed amd64, arm64, or i386
ARCH=$(dpkg --print-architecture)

echo "Detected architecture: $ARCH"

VERSION="5.12.26-stable"

if [[ "$ARCH" == "amd64" ]]; then
  PKG="containDB_${VERSION}_amd64.deb"
elif [[ "$ARCH" == "arm64" ]]; then
  PKG="containDB_${VERSION}_arm64.deb"
elif [[ "$ARCH" == "i386" ]]; then
  PKG="containDB_${VERSION}_i386.deb"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

URL="https://github.com/AnkanSaha/ContainDB/releases/download/${VERSION}/${PKG}"

# Download package
wget -q $URL -O /tmp/$PKG

# Install package
sudo dpkg -i /tmp/$PKG

# Clean up
rm /tmp/$PKG
