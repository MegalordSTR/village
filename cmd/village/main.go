package main

import (
	"fmt"
	"net/http"
	"os"
)

const version = "0.1.0"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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

	// Start health endpoint server
	http.HandleFunc("/health", healthHandler)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Health server error: %v\n", err)
		}
	}()

	// Keep main alive
	select {}
}
