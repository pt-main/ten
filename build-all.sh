#!/bin/bash

cd src

NAME="ten"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT="${NAME}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        OUTPUT+=".exe"
    fi
    echo "Building $OUTPUT..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-X main.version=$VERSION" -o "../bin/$OUTPUT" .
    if [ $? -ne 0 ]; then
        echo "Failed to build $OUTPUT"
        exit 1
    fi
done
echo "Done. Binaries in ./bin/"