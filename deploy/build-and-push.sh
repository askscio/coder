#!/usr/bin/env bash

# Build and push custom Coder image with dev prebuilds support.
#
# Usage:
#   ./build-and-push.sh [IMAGE_REPO] [IMAGE_TAG]
#
# Examples:
#   ./build-and-push.sh                                    # Uses defaults
#   ./build-and-push.sh ghcr.io/rahul-roy-glean/coder      # Custom repo
#   ./build-and-push.sh ghcr.io/rahul-roy-glean/coder dev  # Custom repo and tag

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Defaults
IMAGE_REPO="${1:-ghcr.io/rahul-roy-glean/coder}"
IMAGE_TAG="${2:-dev-prebuilds}"
ARCH="${ARCH:-amd64}"

FULL_IMAGE="${IMAGE_REPO}:${IMAGE_TAG}"

echo "=== Building custom Coder server with dev prebuilds ==="
echo "Repository: ${IMAGE_REPO}"
echo "Tag: ${IMAGE_TAG}"
echo "Architecture: ${ARCH}"
echo ""

cd "$REPO_ROOT"

# Step 1: Build the coder binary
echo "--- Building coder binary for linux/${ARCH}..."
make -j build/coder_linux_${ARCH}

BINARY_PATH="build/coder_linux_${ARCH}"
if [[ ! -f "$BINARY_PATH" ]]; then
    echo "ERROR: Binary not found at $BINARY_PATH"
    exit 1
fi

echo "Binary built: $BINARY_PATH"

# Step 2: Build Docker image
echo ""
echo "--- Building Docker image: ${FULL_IMAGE}..."
./scripts/build_docker.sh \
    --arch "${ARCH}" \
    --target "${FULL_IMAGE}" \
    --build-base "coder-base:local" \
    "$BINARY_PATH"

echo ""
echo "--- Docker image built successfully: ${FULL_IMAGE}"

# Step 3: Push to registry
echo ""
read -p "Push image to registry? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "--- Pushing ${FULL_IMAGE}..."
    docker push "${FULL_IMAGE}"
    echo "--- Image pushed successfully!"
else
    echo "--- Skipping push. You can push manually with:"
    echo "    docker push ${FULL_IMAGE}"
fi

echo ""
echo "=== Build complete ==="
echo ""
echo "To deploy with Helm, update your values file with:"
echo ""
echo "  coder:"
echo "    image:"
echo "      repo: \"${IMAGE_REPO}\""
echo "      tag: \"${IMAGE_TAG}\""
echo "    env:"
echo "      - name: CODER_DEV_PREBUILDS"
echo "        value: \"true\""
echo "      - name: CODER_TELEMETRY_ENABLE"
echo "        value: \"false\""
echo ""
echo "Then run:"
echo "  helm upgrade --install coder ./helm/coder -f deploy/values-dev-prebuilds.yaml"
