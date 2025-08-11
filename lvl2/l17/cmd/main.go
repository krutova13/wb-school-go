package main

import (
	"fmt"
	"os"

	"telnet/internal/parser"
	"telnet/internal/telnet"
)

func main() {
	argParser := parser.NewCommandLineParser()

	config, err := argParser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: telnet [--timeout=10s] host port\n")
		os.Exit(1)
	}

	client := telnet.NewClient(config)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}

	if err := client.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Client error: %v\n", err)
		os.Exit(1)
	}

	if err := client.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing client: %v\n", err)
		os.Exit(1)
	}
}
