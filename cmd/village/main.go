package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Village Simulation v%s\n", version)
		return
	}

	fmt.Println("Village Simulation - Core Engine")
	fmt.Println("=================================")
	fmt.Println("This is the main entry point for the village simulation game.")
	fmt.Println("Development in progress...")
	fmt.Printf("Version: %s\n", version)
}
