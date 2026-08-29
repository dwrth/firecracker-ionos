package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dwrth/knaller/internal/capacity"
	"github.com/dwrth/knaller/internal/config"
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
	case "create", "list", "inspect", "start", "stop", "delete":
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
