#!/usr/bin/env bash
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

die() {
    echo "ERROR: $*" >&2
    exit 1
}

if [[ "$(id -u)" -ne 0 ]]; then
    die "Run this script as root"
fi

echo "==> Checking operating system"

# shellcheck disable=SC1091
source /etc/os-release

[[ "$ID" == "ubuntu" ]] ||
    die "Expected Ubuntu, found: $ID"

if [[ "$VERSION_ID" != "26.04" ]]; then
    echo "WARNING: This build was validated on Ubuntu 26.04."
    echo "Current version: $VERSION_ID"
fi

echo "==> Installing bootstrap/build dependencies"

apt-get update

DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates \
    curl \
    debootstrap \
    e2fsprogs \
    jq \
    iproute2 \
    nftables \
    util-linux \
    kmod \
    xz-utils \
    tar

echo "==> Checking nested virtualization"

"$REPO/host/firecracker-kvm-setup"

[[ -c /dev/kvm ]] ||
    die "/dev/kvm is unavailable"

echo "==> /dev/kvm is available"

echo
echo "==> Fetching Firecracker and jailer"

"$REPO/scripts/fetch-firecracker.sh"

echo
echo "==> Verifying pinned guest kernel"

"$REPO/scripts/verify-kernel.sh"

echo
echo "==> Building fresh guest root filesystems"

"$REPO/scripts/build-guests.sh"

echo
echo "==> Installing Firecracker host"

"$REPO/scripts/install-host.sh"

echo
echo "==> Validating installed host"

"$REPO/scripts/validate-host.sh"

echo
echo "=================================================="
echo "Bootstrap completed successfully."
echo "=================================================="
echo
echo "The VMs are installed and enabled but were not"
echo "explicitly started by this bootstrap."
echo
echo "Recommended next action:"
echo "    reboot"
