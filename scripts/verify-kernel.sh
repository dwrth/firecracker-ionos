#!/usr/bin/env bash
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# shellcheck disable=SC1091
source "$REPO/artifacts/kernel/manifest.env"

KERNEL="$REPO/artifacts/kernel/$KERNEL_FILENAME"

if [[ ! -f "$KERNEL" ]]; then
	echo "Kernel artifact missing: $KERNEL" >&2
	exit 1
fi

ACTUAL="$(sha256sum "$KERNEL" | awk '{print $1}')"

if [[ "$ACTUAL" != "$KERNEL_SHA256" ]]; then
	echo "Kernel checksum mismatch" >&2
	echo "expected: $KERNEL_SHA256" >&2
	echo "actual:   $ACTUAL" >&2
	exit 1
fi

echo "Kernel verified:"
echo "  version: $KERNEL_VERSION"
echo "  file:    $KERNEL_FILENAME"
echo "  sha256:  $ACTUAL"
