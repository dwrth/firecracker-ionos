#!/usr/bin/env bash
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ "$(id -u)" -ne 0 ]]; then
    echo "This installer must run as root" >&2
    exit 1
fi

# shellcheck disable=SC1091
source "$REPO/versions.env"

# shellcheck disable=SC1091
source "$REPO/artifacts/kernel/manifest.env"

FC_BIN="$REPO/build/artifacts/firecracker"
JAILER_BIN="$REPO/build/artifacts/jailer"
KERNEL="$REPO/artifacts/kernel/$KERNEL_FILENAME"

VM1_IMAGE="$REPO/build/guests/vm1.ext4"
VM2_IMAGE="$REPO/build/guests/vm2.ext4"

VM1_CONFIG="$REPO/guests/vm1/config.json"
VM2_CONFIG="$REPO/guests/vm2/config.json"

die() {
    echo "ERROR: $*" >&2
    exit 1
}

verify_file() {
    [[ -f "$1" ]] || die "Required file missing: $1"
}

verify_artifacts() {
    echo "==> Verifying required artifacts"

    verify_file "$FC_BIN"
    verify_file "$JAILER_BIN"
    verify_file "$KERNEL"

    verify_file "$VM1_IMAGE"
    verify_file "$VM2_IMAGE"

    verify_file "$VM1_CONFIG"
    verify_file "$VM2_CONFIG"

    "$REPO/scripts/verify-kernel.sh"

    EXPECTED_FC="$(
        awk '$1 == "firecracker" {print $2}' \
            "$REPO/artifacts/artifacts.sha256"
    )"

    EXPECTED_JAILER="$(
        awk '$1 == "jailer" {print $2}' \
            "$REPO/artifacts/artifacts.sha256"
    )"

    ACTUAL_FC="$(sha256sum "$FC_BIN" | awk '{print $1}')"
    ACTUAL_JAILER="$(sha256sum "$JAILER_BIN" | awk '{print $1}')"

    [[ "$ACTUAL_FC" == "$EXPECTED_FC" ]] ||
        die "Firecracker checksum mismatch"

    [[ "$ACTUAL_JAILER" == "$EXPECTED_JAILER" ]] ||
        die "Jailer checksum mismatch"

    (
        cd "$REPO/build/guests"
        sha256sum -c "$REPO/artifacts/guest-images.sha256"
    )

    jq empty "$VM1_CONFIG"
    jq empty "$VM2_CONFIG"

    echo "==> Artifact verification passed"
}

install_packages() {
    echo "==> Installing host packages"

    apt-get update

    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        ca-certificates \
        curl \
        jq \
        iproute2 \
        nftables \
        e2fsprogs \
        util-linux \
        kmod
}

ensure_account() {
    NAME="$1"
    UID_NUM="$2"
    GID_NUM="$3"

    if getent group "$NAME" >/dev/null; then
        CURRENT_GID="$(getent group "$NAME" | cut -d: -f3)"

        [[ "$CURRENT_GID" == "$GID_NUM" ]] ||
            die "$NAME group exists with GID $CURRENT_GID, expected $GID_NUM"
    elif getent group "$GID_NUM" >/dev/null; then
        die "GID $GID_NUM already belongs to another group"
    else
        groupadd \
            --system \
            --gid "$GID_NUM" \
            "$NAME"
    fi

    if getent passwd "$NAME" >/dev/null; then
        CURRENT_UID="$(id -u "$NAME")"

        [[ "$CURRENT_UID" == "$UID_NUM" ]] ||
            die "$NAME user exists with UID $CURRENT_UID, expected $UID_NUM"
    elif getent passwd "$UID_NUM" >/dev/null; then
        die "UID $UID_NUM already belongs to another user"
    else
        useradd \
            --system \
            --uid "$UID_NUM" \
            --gid "$GID_NUM" \
            --no-create-home \
            --shell /usr/sbin/nologin \
            "$NAME"
    fi
}

install_accounts() {
    echo "==> Creating Firecracker identities"

    ensure_account fc-vm1 "$VM1_UID" "$VM1_GID"
    ensure_account fc-vm2 "$VM2_UID" "$VM2_GID"
}

install_binaries() {
    echo "==> Installing Firecracker binaries"

    install \
        -o root \
        -g root \
        -m 0755 \
        "$FC_BIN" \
        /usr/local/bin/firecracker

    install \
        -o root \
        -g root \
        -m 0755 \
        "$JAILER_BIN" \
        /usr/local/bin/jailer

    /usr/local/bin/firecracker --version
    /usr/local/bin/jailer --version
}

install_host_helpers() {
    echo "==> Installing host helpers"

    for SCRIPT in \
        firecracker-kvm-setup \
        firecracker-host-network \
        firecracker-vm-network \
        firecracker-prepare \
        firecracker-start \
        firecracker-stop
    do
        install \
            -o root \
            -g root \
            -m 0755 \
            "$REPO/host/$SCRIPT" \
            "/usr/local/sbin/$SCRIPT"
    done
}

install_host_configuration() {
    echo "==> Installing kernel/module/sysctl configuration"

    install \
        -o root \
        -g root \
        -m 0644 \
        "$REPO/host/modules-load/firecracker-kvm.conf" \
        /etc/modules-load.d/firecracker-kvm.conf

    install \
        -o root \
        -g root \
        -m 0644 \
        "$REPO/host/sysctl/99-firecracker.conf" \
        /etc/sysctl.d/99-firecracker.conf

    sysctl --system >/dev/null

    /usr/local/sbin/firecracker-kvm-setup

    [[ -c /dev/kvm ]] ||
        die "/dev/kvm is unavailable"
}

install_vm() {
    VM="$1"
    UID_NUM="$2"
    GID_NUM="$3"
    IMAGE="$4"
    CONFIG="$5"

    DEST="/var/lib/firecracker/$VM"

    echo "==> Installing $VM"

    install -d \
        -o root \
        -g root \
        -m 0750 \
        "$DEST"

    cp --reflink=auto --sparse=always \
        "$IMAGE" \
        "$DEST/rootfs.ext4"

    install \
        -o root \
        -g "$GID_NUM" \
        -m 0440 \
        "$KERNEL" \
        "$DEST/vmlinux"

    install \
        -o root \
        -g "$GID_NUM" \
        -m 0440 \
        "$CONFIG" \
        "$DEST/config.json"

    chown "$UID_NUM:$GID_NUM" \
        "$DEST/rootfs.ext4"

    chmod 0600 \
        "$DEST/rootfs.ext4"
}

install_vm_artifacts() {
    echo "==> Installing VM artifacts"

    install -d \
        -o root \
        -g root \
        -m 0755 \
        /var/lib/firecracker

    install -d \
        -o root \
        -g root \
        -m 0755 \
        /srv/jailer

    install -d \
        -o root \
        -g root \
        -m 0755 \
        /srv/jailer/firecracker

    # firecracker-prepare uses hard links.
    VAR_DEV="$(stat -c %d /var/lib/firecracker)"
    JAIL_DEV="$(stat -c %d /srv/jailer/firecracker)"

    [[ "$VAR_DEV" == "$JAIL_DEV" ]] ||
        die "/var/lib/firecracker and /srv/jailer must be on the same filesystem"

    install_vm \
        vm1 \
        "$VM1_UID" \
        "$VM1_GID" \
        "$VM1_IMAGE" \
        "$VM1_CONFIG"

    install_vm \
        vm2 \
        "$VM2_UID" \
        "$VM2_GID" \
        "$VM2_IMAGE" \
        "$VM2_CONFIG"
}

install_systemd() {
    echo "==> Installing systemd units"

    install \
        -o root \
        -g root \
        -m 0644 \
        "$REPO/host/systemd/firecracker-host-network.service" \
        /etc/systemd/system/firecracker-host-network.service

    install \
        -o root \
        -g root \
        -m 0644 \
        "$REPO/host/systemd/firecracker-network@.service" \
        '/etc/systemd/system/firecracker-network@.service'

    install \
        -o root \
        -g root \
        -m 0644 \
        "$REPO/host/systemd/firecracker@.service" \
        '/etc/systemd/system/firecracker@.service'

    # Remove the obsolete pre-template network unit if present.
    systemctl disable firecracker-network.service \
        >/dev/null 2>&1 || true

    rm -f /etc/systemd/system/firecracker-network.service

    systemctl daemon-reload

    systemctl enable firecracker-host-network.service
    systemctl enable firecracker@vm1.service
    systemctl enable firecracker@vm2.service
}

validate_installation() {
    echo "==> Validating installation"

    test -x /usr/local/bin/firecracker
    test -x /usr/local/bin/jailer

    test -x /usr/local/sbin/firecracker-host-network
    test -x /usr/local/sbin/firecracker-vm-network
    test -x /usr/local/sbin/firecracker-prepare
    test -x /usr/local/sbin/firecracker-start
    test -x /usr/local/sbin/firecracker-stop

    test -c /dev/kvm

    test -f /var/lib/firecracker/vm1/vmlinux
    test -f /var/lib/firecracker/vm1/rootfs.ext4
    test -f /var/lib/firecracker/vm1/config.json

    test -f /var/lib/firecracker/vm2/vmlinux
    test -f /var/lib/firecracker/vm2/rootfs.ext4
    test -f /var/lib/firecracker/vm2/config.json

    jq empty /var/lib/firecracker/vm1/config.json
    jq empty /var/lib/firecracker/vm2/config.json

    systemctl is-enabled firecracker-host-network.service
    systemctl is-enabled firecracker@vm1.service
    systemctl is-enabled firecracker@vm2.service

    echo "==> Installation validation passed"
}

main() {
    verify_artifacts
    install_packages
    install_accounts
    install_binaries
    install_host_helpers
    install_host_configuration
    install_vm_artifacts
    install_systemd
    validate_installation

    echo
    echo "Host installation complete."
    echo
    echo "The Firecracker VMs have been installed and enabled,"
    echo "but this installer intentionally does not start them."
    echo
    echo "To test:"
    echo "  systemctl start firecracker-host-network.service"
    echo "  systemctl start firecracker@vm1.service"
    echo "  systemctl start firecracker@vm2.service"
}

main "$@"
