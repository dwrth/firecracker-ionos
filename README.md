# Firecracker on IONOS Cloud Cubes

Reproducible host configuration for running two jailed Firecracker
microVMs on an IONOS Cloud Cube.

## Layout

- `host/` - host networking, jailer helpers and systemd configuration
- `guests/` - per-VM Firecracker configurations
- `scripts/` - image/build automation
- `artifacts/` - artifact manifests and checksums
- `versions.env` - pinned software and VM versions

## Runtime layout

VM1:
- UID/GID: 1101
- network namespace: fc-vm1
- guest address: 172.16.1.2/30
- gateway: 172.16.1.1

VM2:
- UID/GID: 1102
- network namespace: fc-vm2
- guest address: 172.16.2.2/30
- gateway: 172.16.2.1

Both VMs are started through Firecracker jailer and systemd.
