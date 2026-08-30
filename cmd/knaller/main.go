package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/capacity"
	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/scheduler"
	"github.com/dwrth/knaller/internal/state"
)

const defaultConfigPath = "/etc/knaller/config.yaml"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return 0
	case "capacity":
		return runCapacity(args[1:])
	case "config":
		return runConfig(args[1:])
	case "create":
		return runCreate(args[1:])
	case "list", "inspect", "start", "stop", "delete":
		fmt.Fprintf(os.Stderr, "knaller %s: not implemented yet\n", args[0])
		return 1
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage(os.Stderr)
		return 2
	}
}

func runCapacity(args []string) int {
	fs := flag.NewFlagSet("capacity", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: knaller capacity [flags]
Flags:
`)
		fs.PrintDefaults()
	}
	configPath := fs.String("config", defaultConfigPath, "path to knaller config file")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	var cfg config.Config
	if err := cfg.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		return 1
	}

	cap, err := capacity.Collect(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to collect capacity: %v\n", err)
		return 1
	}

	capacity.Print(os.Stdout, cap)
	return 0
}

func runConfig(args []string) int {
	if len(args) == 0 {
		printConfigUsage(os.Stderr)
		return 2
	}

	switch args[0] {
	case "validate":
		return runConfigValidate(args[1:])
	case "help", "-h", "--help":
		printConfigUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printConfigUsage(os.Stderr)
		return 2
	}
}

func runConfigValidate(args []string) int {
	fs := flag.NewFlagSet("config validate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: knaller config validate [flags]
Flags:
`)
		fs.PrintDefaults()
	}
	configPath := fs.String("config", defaultConfigPath, "path to knaller config file")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	var cfg config.Config
	if err := cfg.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		return 1
	}

	fmt.Fprintln(os.Stdout, "configuration is valid")
	return 0
}

func runCreate(args []string) int {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: knaller create [flags]
Flags:
`)
		fs.PrintDefaults()
	}
	configPath := fs.String("config", defaultConfigPath, "path to knaller config file")
	name := fs.String("name", "", "VM name")
	cpus := fs.Int("cpus", 0, "vCPUs")
	memory := fs.Int("memory", 0, "memory in MiB")
	dryRun := fs.Bool("dry-run", false, "print allocation without writing state")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if *name == "" || *cpus == 0 || *memory == 0 {
		fmt.Fprintln(os.Stderr, "create: --name, --cpus, and --memory are required")
		fs.Usage()
		return 2
	}
	if !*dryRun {
		fmt.Fprintln(os.Stderr, "knaller create: only --dry-run is supported for now")
		return 1
	}

	var cfg config.Config
	if err := cfg.Load(*configPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		return 1
	}

	store := state.New(cfg.State.Directory)
	existing, err := store.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list VMs: %v\n", err)
		return 1
	}

	cap, err := capacity.Collect(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to collect capacity: %v\n", err)
		return 1
	}

	if err := scheduler.Admit(&cfg, cap, scheduler.Request{VCPUs: *cpus, MemoryMiB: *memory}); err != nil {
		fmt.Fprintf(os.Stderr, "admission denied: %v\n", err)
		return 1
	}

	vm, err := allocate.Build(existing, &cfg, *name, *cpus, *memory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "allocation failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "dry-run: would create %s (%s)\n", vm.ID, vm.Name)
	fmt.Fprintf(os.Stdout, "  uid/gid        %d/%d\n", vm.UID, vm.GID)
	fmt.Fprintf(os.Stdout, "  vcpus/memory   %d / %d MiB\n", vm.VCPUs, vm.MemoryMiB)
	fmt.Fprintf(os.Stdout, "  guest          %s  gw %s  subnet %s\n", vm.GuestIP, vm.GatewayIP, vm.GuestSubnet)
	fmt.Fprintf(os.Stdout, "  transit        host %s  ns %s  subnet %s\n", vm.TransitHostIP, vm.TransitNSIP, vm.TransitSubnet)
	fmt.Fprintf(os.Stdout, "  netns/veth/tap %s  %s/%s  %s\n", vm.Namespace, vm.HostVeth, vm.NSVeth, vm.TAP)
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: knaller <command> [flags]
Commands:
  capacity   Show host and VM resource capacity
  config     Manage configuration
  create     Create a new microVM
  list       List microVMs
  inspect    Show details for a microVM
  start      Start a microVM
  stop       Stop a microVM
  delete     Delete a microVM
`)
}

func printConfigUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: knaller config <command> [flags]
Commands:
  validate   Validate the configuration
`)
}
