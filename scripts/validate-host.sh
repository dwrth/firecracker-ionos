#!/usr/bin/env bash
set -euo pipefail

fail=0

check() {
    DESCRIPTION="$1"
    shift

    printf '%-45s' "$DESCRIPTION"

    if "$@" >/dev/null 2>&1; then
        echo "OK"
    else
        echo "FAIL"
        fail=1
    fi
}

echo "=== FIRECRACKER HOST VALIDATION ==="
echo

check \
    "/dev/kvm available" \
    test -c /dev/kvm

check \
    "Firecracker installed" \
    test -x /usr/local/bin/firecracker

check \
    "Jailer installed" \
    test -x /usr/local/bin/jailer

check \
    "VM1 rootfs installed" \
    test -f /var/lib/firecracker/vm1/rootfs.ext4

check \
    "VM2 rootfs installed" \
    test -f /var/lib/firecracker/vm2/rootfs.ext4

check \
    "VM1 kernel installed" \
    test -f /var/lib/firecracker/vm1/vmlinux

check \
    "VM2 kernel installed" \
    test -f /var/lib/firecracker/vm2/vmlinux

check \
    "VM1 config valid" \
    jq empty /var/lib/firecracker/vm1/config.json

check \
    "VM2 config valid" \
    jq empty /var/lib/firecracker/vm2/config.json

check \
    "Host network enabled" \
    systemctl is-enabled firecracker-host-network.service

check \
    "VM1 enabled" \
    systemctl is-enabled firecracker@vm1.service

check \
    "VM2 enabled" \
    systemctl is-enabled firecracker@vm2.service

echo

EXPECTED_KERNEL="$(
    awk -F= \
        '$1 == "KERNEL_SHA256" {print $2}' \
        "$(dirname "$0")/../artifacts/kernel/manifest.env"
)"

VM1_KERNEL="$(
    sha256sum /var/lib/firecracker/vm1/vmlinux |
    awk '{print $1}'
)"

VM2_KERNEL="$(
    sha256sum /var/lib/firecracker/vm2/vmlinux |
    awk '{print $1}'
)"

printf '%-45s' "VM1 kernel checksum"
if [[ "$VM1_KERNEL" == "$EXPECTED_KERNEL" ]]; then
    echo "OK"
else
    echo "FAIL"
    fail=1
fi

printf '%-45s' "VM2 kernel checksum"
if [[ "$VM2_KERNEL" == "$EXPECTED_KERNEL" ]]; then
    echo "OK"
else
    echo "FAIL"
    fail=1
fi

echo

if [[ "$fail" -ne 0 ]]; then
    echo "Host validation FAILED"
    exit 1
fi

echo "Host validation PASSED"
