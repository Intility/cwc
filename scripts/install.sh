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

URL="https://github.com/emilkje/cwc/releases/latest/download/cwc-$OS-$ARCH$EXT"

# Download the correct binary
wget "$URL" -O cwc-$OS-$ARCH$EXT

# Make it executable (not necessary for Windows)
if [ "$OS" != "windows" ]; then
    chmod +x cwc-$OS-$ARCH$EXT
fi

echo "Downloaded cwc for $OS-$ARCH"
