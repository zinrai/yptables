package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zinrai/yptables/internal/config"
	"github.com/zinrai/yptables/internal/generator"
)

func main() {
	// Command line flags
	formatFlag := flag.String("format", "script", "Output format: 'script' or 'restore'")
	outputFlag := flag.String("output", "", "Output file (default: stdout)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <config.yaml>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Determine output format
	var format generator.Format
	switch *formatFlag {
	case "script":
		format = generator.ShellScript
	case "restore":
		format = generator.IPTablesRestore
	default:
		log.Fatalf("Invalid format: %s", *formatFlag)
	}

	// Load configuration
	cfg, err := config.LoadFromFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Generate commands
	gen := generator.New(format)
	commands, err := gen.Generate(cfg)
	if err != nil {
		log.Fatalf("Failed to generate commands: %v", err)
	}

	// Output the result
	content := strings.Join(commands, "\n") + "\n"
	if *outputFlag == "" {
		fmt.Print(content)
	} else {
		if err := os.WriteFile(*outputFlag, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}
}
