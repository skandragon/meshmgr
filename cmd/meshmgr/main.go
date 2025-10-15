package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Meshtastic Node Manager")
	fmt.Println("Version: 0.1.0")

	if len(os.Args) > 1 {
		fmt.Printf("Command: %s\n", os.Args[1])
		fmt.Println("(Command handling to be implemented)")
	} else {
		fmt.Println("Usage: meshmgr <command>")
		fmt.Println("Commands will be implemented soon.")
	}
}
