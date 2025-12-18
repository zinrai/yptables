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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <config.yaml>\n", os.Args[0])
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadFromFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Generate iptables-restore format
	gen := generator.New()
	lines, err := gen.Generate(cfg)
	if err != nil {
		log.Fatalf("Failed to generate iptables-restore output: %v", err)
	}

	// Output to stdout
	content := strings.Join(lines, "\n") + "\n"
	fmt.Print(content)
}
