package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Automation Platform CLI (atp)")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  atp <command>")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  status    Show platform status")
		fmt.Println("  version   Show version")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("atp v0.1.0")
	case "status":
		fmt.Println("Use the API at /api/v1 or the dashboard at http://localhost:3000")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
