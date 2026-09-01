#!/usr/bin/env bash
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD="$REPO/build"
BASE="$BUILD/rootfs-base"
OUTPUT="$BUILD/guests"

UBUNTU_RELEASE="noble"
UBUNTU_MIRROR="http://archive.ubuntu.com/ubuntu"

if [[ "$(uname -m)" != "x86_64" ]]; then
	echo "This build currently expects x86_64" >&2
	exit 1
fi

echo "==> Cleaning previous build"

rm -rf "$BASE" "$OUTPUT"

mkdir -p \
	"$BASE" \
	"$OUTPUT"

echo "==> Bootstrapping Ubuntu ${UBUNTU_RELEASE}"

debootstrap \
	--arch=amd64 \
	--variant=minbase \
	--include=systemd-sysv,udev,iproute2,iputils-ping,ca-certificates,curl,procps \
	"$UBUNTU_RELEASE" \
	"$BASE" \
	"$UBUNTU_MIRROR"

echo "==> Configuring Ubuntu repositories"

cat >"$BASE/etc/apt/sources.list" <<'APT'
deb http://archive.ubuntu.com/ubuntu noble main universe
deb http://archive.ubuntu.com/ubuntu noble-updates main universe
deb http://security.ubuntu.com/ubuntu noble-security main universe
APT

# Give the build chroot working DNS.
rm -f "$BASE/etc/resolv.conf"

cat >"$BASE/etc/resolv.conf" <<'RESOLV'
nameserver 1.1.1.1
nameserver 8.8.8.8
RESOLV

echo "==> Updating guest packages"

chroot "$BASE" /usr/bin/env \
	DEBIAN_FRONTEND=noninteractive \
	apt-get update

chroot "$BASE" /usr/bin/env \
	DEBIAN_FRONTEND=noninteractive \
	apt-get -y dist-upgrade

echo "==> Configuring common guest settings"

mkdir -p "$BASE/etc/systemd/network"
mkdir -p "$BASE/var/lib/dbus"

# Each instantiated guest must generate its own machine-id.
rm -f "$BASE/etc/machine-id"
touch "$BASE/etc/machine-id"

rm -f "$BASE/var/lib/dbus/machine-id"
ln -s /etc/machine-id "$BASE/var/lib/dbus/machine-id"

# Do not bake SSH host identity into the image.
rm -f "$BASE/etc/ssh/ssh_host_"* 2>/dev/null || true

# Use static resolvers for this minimal image.
rm -f "$BASE/etc/resolv.conf"

cat >"$BASE/etc/resolv.conf" <<'RESOLV'
nameserver 1.1.1.1
nameserver 8.8.8.8
options timeout:2 attempts:2
RESOLV

if [[ ! -x "$BASE/usr/lib/systemd/systemd-udevd" ]]; then
	echo "systemd-udevd is missing from guest image" >&2
	exit 1
fi

# Networking will be handled by systemd-networkd.
systemctl --root="$BASE" enable systemd-networkd.service

# Useful serial console for Firecracker diagnostics.
systemctl --root="$BASE" enable serial-getty@ttyS0.service

# Do not bake a known root password into the image.
chroot "$BASE" passwd -l root >/dev/null

echo "==> Recording installed package versions"

dpkg-query \
	--admindir="$BASE/var/lib/dpkg" \
	-W \
	-f='${Package}\t${Version}\n' |
	sort \
		>"$REPO/artifacts/guest-packages.tsv"

echo "==> Cleaning base image"

chroot "$BASE" apt-get clean

rm -rf "$BASE/var/lib/apt/lists/"*
rm -rf "$BASE/var/cache/apt/"*
rm -rf "$BASE/tmp/"*
rm -rf "$BASE/var/tmp/"*
rm -rf "$BASE/var/log/"*

mkdir -p \
	"$BASE/var/log" \
	"$BASE/tmp" \
	"$BASE/var/tmp"

chmod 1777 \
	"$BASE/tmp" \
	"$BASE/var/tmp"

echo "==> Creating base rootfs image"

IMAGE="$OUTPUT/base-rootfs.ext4"

truncate -s 1G "$IMAGE"

mkfs.ext4 \
	-q \
	-F \
	-L knaller-base \
	-d "$BASE" \
	"$IMAGE"

e2fsck -fn "$IMAGE" >/dev/null

echo "    created $IMAGE"

echo "==> Recording image hash"

(
	cd "$OUTPUT"
	sha256sum base-rootfs.ext4
) >"$REPO/artifacts/guest-images.sha256"
echo
echo "Build complete:"
ls -lh "$OUTPUT"
echo
echo "SHA256:"
cat "$REPO/artifacts/guest-images.sha256"
