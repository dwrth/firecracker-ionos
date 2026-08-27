package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dwrth/knaller/internal/capacity"
)

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
		capacity.PrintStub()
		return 0
	case "create", "list", "inspect", "start", "stop", "delete":
		fmt.Fprintf(os.Stderr, "knaller %s: not implemented yet\n", args[0])
		return 1
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage(os.Stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: knaller <command>
Commands:
  capacity   Show host and VM resource capacity
  create     Create a new microVM
  list       List microVMs
  inspect    Show details for a microVM
  start      Start a microVM
  stop       Stop a microVM
  delete     Delete a microVM
`)
}
