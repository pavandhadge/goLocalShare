#!/bin/bash

# Build script for cross-compiling a Go application for multiple platforms

# Application information
APP_NAME="localshare"  # Change this to your application name
VERSION="1.0.0"
OUTPUT_DIR="builds"

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# List of target platforms and architectures
PLATFORMS=(
    "linux/amd64"
    "linux/386"
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    "darwin/arm64"  # For M1/M2 Macs
)

# Function to build for a specific platform
build_for_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME="$APP_NAME-$VERSION-$GOOS-$GOARCH"

    # Special handling for Windows
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi

    echo "Building for $GOOS/$GOARCH..."

    env GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME" main.go

    # Compress the binary (optional)
    if command -v upx &> /dev/null; then
        echo "Compressing binary with UPX..."
        upx --best "$OUTPUT_DIR/$OUTPUT_NAME"
    fi
}

# Build for all platforms
for platform in "${PLATFORMS[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    build_for_platform $GOOS $GOARCH
done

# Create checksums
echo "Creating checksums..."
cd "$OUTPUT_DIR" && sha256sum * > checksums.sha256 && cd ..

echo "Build complete! Files are in the $OUTPUT_DIR directory:"
ls -lh "$OUTPUT_DIR"
