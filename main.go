package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

	cfg, err := LoadFromFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	gen := NewGenerator()
	lines, err := gen.Generate(cfg)
	if err != nil {
		log.Fatalf("Failed to generate iptables-restore output: %v", err)
	}

	content := strings.Join(lines, "\n") + "\n"
	fmt.Print(content)
}
