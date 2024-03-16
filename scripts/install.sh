#!/bin/bash

# Detect the operating system and architecture
OS="unknown"
ARCH=$(uname -m)

case "$(uname -s)" in
    Darwin) OS="darwin";;
    Linux) OS="linux";;
    CYGWIN*|MINGW32*|MSYS*|MINGW*) OS="windows";;
esac

# Map architecture names to those used in the releases
case "$ARCH" in
    x86_64) ARCH="amd64";;
    arm64) ARCH="arm64";;
esac

# For Windows, change the executable name
EXT=""
if [ "$OS" = "windows" ]; then
    EXT=".exe"
fi

URL="https://github.com/intility/cwc/releases/latest/download/cwc-$OS-$ARCH$EXT"

# Download the correct binary
if command -v wget > /dev/null; then
    wget "$URL" -O "cwc$EXT"
elif command -v curl > /dev/null; then
    curl -L "$URL" -o "cwc$EXT"
else
    echo "Error: Neither wget nor curl is installed."
    exit 1
fi

# Make it executable (not necessary for Windows)
if [ "$OS" != "windows" ]; then
    chmod +x "cwc$EXT"
fi

echo "Downloaded cwc for $OS-$ARCH"
