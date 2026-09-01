#!/usr/bin/env bash
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# shellcheck disable=SC1091
source "$REPO/versions.env"

ARCH="$(uname -m)"

if [[ "$ARCH" != "x86_64" ]]; then
	echo "Unsupported architecture: $ARCH" >&2
	exit 1
fi

VERSION="$FIRECRACKER_VERSION"
ARCHIVE="firecracker-${VERSION}-${ARCH}.tgz"

BASE_URL="https://github.com/firecracker-microvm/firecracker/releases/download/${VERSION}"

DOWNLOAD_DIR="$REPO/build/downloads/firecracker"
OUTPUT_DIR="$REPO/build/artifacts"

mkdir -p "$DOWNLOAD_DIR" "$OUTPUT_DIR"

echo "==> Firecracker version: $VERSION"
echo "==> Architecture: $ARCH"

echo "==> Downloading release archive"

curl -fL \
	"${BASE_URL}/${ARCHIVE}" \
	-o "$DOWNLOAD_DIR/$ARCHIVE"

echo "==> Downloading official archive checksum"

curl -fL \
	"${BASE_URL}/${ARCHIVE}.sha256.txt" \
	-o "$DOWNLOAD_DIR/${ARCHIVE}.sha256.txt"

echo "==> Verifying official release archive"

(
	cd "$DOWNLOAD_DIR"
	sha256sum -c "${ARCHIVE}.sha256.txt"
)

echo "==> Extracting"

rm -rf "$DOWNLOAD_DIR/extracted"
mkdir -p "$DOWNLOAD_DIR/extracted"

tar -xzf "$DOWNLOAD_DIR/$ARCHIVE" \
	-C "$DOWNLOAD_DIR/extracted"

FC_SOURCE="$(
	find "$DOWNLOAD_DIR/extracted" \
		-type f \
		-name "firecracker-${VERSION}-${ARCH}" \
		-print \
		-quit
)"

JAILER_SOURCE="$(
	find "$DOWNLOAD_DIR/extracted" \
		-type f \
		-name "jailer-${VERSION}-${ARCH}" \
		-print \
		-quit
)"

if [[ -z "$FC_SOURCE" || ! -f "$FC_SOURCE" ]]; then
	echo "Unable to locate Firecracker binary" >&2
	exit 1
fi

if [[ -z "$JAILER_SOURCE" || ! -f "$JAILER_SOURCE" ]]; then
	echo "Unable to locate jailer binary" >&2
	exit 1
fi

install -m 0755 \
	"$FC_SOURCE" \
	"$OUTPUT_DIR/firecracker"

install -m 0755 \
	"$JAILER_SOURCE" \
	"$OUTPUT_DIR/jailer"

echo "==> Verifying binaries against known-good hashes"

EXPECTED_FC="$(
	awk '$1 == "firecracker" {print $2}' \
		"$REPO/artifacts/artifacts.sha256"
)"

EXPECTED_JAILER="$(
	awk '$1 == "jailer" {print $2}' \
		"$REPO/artifacts/artifacts.sha256"
)"

if [[ -z "$EXPECTED_FC" || -z "$EXPECTED_JAILER" ]]; then
	echo "Missing expected hashes in artifacts/artifacts.sha256" >&2
	exit 1
fi

ACTUAL_FC="$(sha256sum "$OUTPUT_DIR/firecracker" | awk '{print $1}')"
ACTUAL_JAILER="$(sha256sum "$OUTPUT_DIR/jailer" | awk '{print $1}')"

if [[ "$ACTUAL_FC" != "$EXPECTED_FC" ]]; then
	echo "Firecracker binary hash mismatch" >&2
	echo "expected: $EXPECTED_FC" >&2
	echo "actual:   $ACTUAL_FC" >&2
	exit 1
fi

if [[ "$ACTUAL_JAILER" != "$EXPECTED_JAILER" ]]; then
	echo "Jailer binary hash mismatch" >&2
	echo "expected: $EXPECTED_JAILER" >&2
	echo "actual:   $ACTUAL_JAILER" >&2
	exit 1
fi

echo "==> Checking versions"

"$OUTPUT_DIR/firecracker" --version
"$OUTPUT_DIR/jailer" --version

echo
echo "==> Verified artifacts"
sha256sum \
	"$OUTPUT_DIR/firecracker" \
	"$OUTPUT_DIR/jailer"
