package main

import (
	"fmt"
	"os"

	"github.com/aptible/mobs/src/service"
	"github.com/aptible/mobs/src/storage/cloverdb"
)

func main() {
	if err := os.MkdirAll("./data", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create data directory: %v\n", err)
		os.Exit(1)
	}
	store, err := cloverdb.NewStore("./data/mobs.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}
	service := service.NewTenantService(store)

	rootCmd := setupCommands(service)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
