package main

import (
	"fmt"
	"os"

	"github.com/aptible/mobs/src/service"
	"github.com/aptible/mobs/src/storage/cloverdb"
)

// main is the entry point for the Multi-tenant Object Storage CLI.
func main() {
	exitOnErr := func(msg string, err error) {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}

	if err := os.MkdirAll("./data", 0755); err != nil {
		exitOnErr("Failed to create data directory", err)
	}
	store, err := cloverdb.NewStore("./data/mobs.db")
	if err != nil {
		exitOnErr("Failed to initialize storage", err)
	}
	service := service.NewTenantService(store)

	rootCmd := setupCommands(service)
	if err := rootCmd.Execute(); err != nil {
		exitOnErr("Error", err)
	}
}
