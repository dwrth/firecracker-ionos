# Knaller

Dynamic Firecracker microVM manager for an IONOS Cloud Cube.

> Knaller is the German word for "firecracker" (/ˈknalɐ/, roughly _KNALL-er_).

A single Ubuntu host runs many isolated microVMs. The `knaller` CLI
creates, sizes, starts, stops, lists, inspects, and deletes VMs based on
host capacity. Each VM runs under Firecracker jailer with its own Linux
user, network namespace, cloned rootfs, and TAP/veth networking. systemd
supervises the lifecycle via `knaller@…` units.

The repo evolved from a static two-VM setup (`v1.0.0`) into this generic
manager. See the roadmap for current progress.

## CLI

```text
knaller capacity          # host resources and what is still free
knaller create            # provision a new VM
knaller list              # running and stopped VMs
knaller inspect <name>    # full VM state and paths
knaller start <name>      # start via systemd
knaller stop <name>       # graceful shutdown
knaller delete <name>     # tear down VM and storage
knaller config validate   # check /etc/knaller/config.yaml
```

Creation examples (target):

```bash
knaller create --name web-1 --cpus 2 --memory 512
knaller create --name worker-1 --size small
knaller create --name worker-2 --auto
```

Admission enforces strict RAM limits and configurable CPU overcommit.
Provisioning clones a shared base rootfs, assigns IDs, UID/GID, guest and
transit subnets, generates Firecracker config, wires networking and NAT,
and persists state for reboot-safe operation.

## Repository layout

- `cmd/knaller/` - CLI entrypoint
- `internal/` - config, capacity, scheduler, state, allocate, storage, ...
- `config/` - example host configuration
- `host/` - jailer helpers, host networking, systemd units
- `guests/` - legacy static VM configs (`vm1`, `vm2`; replaced by dynamic flow)
- `scripts/` - image and host build automation
- `artifacts/` - artifact manifests and checksums
- `versions.env` - pinned software and kernel versions

## Host runtime layout

```text
/etc/knaller/config.yaml

/var/lib/knaller/images/
  base-rootfs.ext4
  vmlinux-6.18.44

/var/lib/knaller/state/
  vm-0001.json
  vm-0002.json

/var/lib/knaller/vms/<vm-id>/
  config.json
  rootfs.ext4
  vmlinux

/srv/jailer/firecracker/<vm-id>/   # disposable jail state

systemd: knaller@.service, knaller-network@.service
```

Each VM gets a deterministic `vm-NNNN` ID, dedicated UID/GID, network
namespace (`kn-vm-NNNN`), guest `/30` on `172.16.x.x`, and transit
`/30` on `10.200.x.x` for host↔namespace veth plumbing.

## Roadmap

- [x] v1.0.0 static two-VM baseline
- [x] Go CLI skeleton
- [x] Config loading
- [x] Capacity reporting
- [x] VM state model
- [x] Scheduler / admission
- [x] Resource allocation
- [x] Base-rootfs workflow
- [x] Firecracker config
- [ ] Jailer scripts
- [ ] Networking / NAT
- [ ] create
- [ ] list / inspect
- [ ] start / stop
- [ ] delete
- [ ] Named sizes
- [ ] Auto sizing
- [ ] CPU controls
- [ ] doctor / reconcile
- [ ] Golden image / clean-room test
