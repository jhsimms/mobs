package main

import (
	"fmt"
	"os"

	"github.com/aptible/mobs/src/service"
	"github.com/spf13/cobra"
)

// setupCommands configures the root and tenant commands for the CLI.
func setupCommands(tenantService *service.TenantService) *cobra.Command {
	exitOnErr := func(msg string, err error) {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "mobs",
		Short: "Multi-tenant Object Storage CLI",
	}

	// Tenant command
	tenantCmd := &cobra.Command{
		Use:   "tenant",
		Short: "Manage tenants",
	}

	// Create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new tenant",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Fprintln(os.Stderr, "--name is required")
				os.Exit(1)
			}
			meta, err := tenantService.CreateTenant(name)
			if err != nil {
				exitOnErr("Failed to create tenant", err)
			}
			fmt.Printf("Tenant created: ID=%s, Name=%s\n", meta.ID, meta.Name)
		},
	}
	createCmd.Flags().String("name", "", "Tenant name")

	// List
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tenants",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := tenantService.ListTenants()
			if err != nil {
				exitOnErr("Failed to list tenants", err)
			}
			if len(list) == 0 {
				fmt.Println("No tenants found.")
				return
			}
			for _, t := range list {
				fmt.Printf("ID=%s, Name=%s\n", t.ID, t.Name)
			}
		},
	}

	// Get
	getCmd := &cobra.Command{
		Use:   "get [TENANT_ID]",
		Short: "Get tenant details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			meta, err := tenantService.GetTenant(id)
			if err != nil {
				exitOnErr("Failed to get tenant", err)
			}
			fmt.Printf("ID=%s, Name=%s\n", meta.ID, meta.Name)
		},
	}

	// Delete
	deleteCmd := &cobra.Command{
		Use:   "delete [TENANT_ID]",
		Short: "Delete a tenant",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			err := tenantService.DeleteTenant(id)
			if err != nil {
				exitOnErr("Failed to delete tenant", err)
			}
			fmt.Printf("Tenant deleted: %s\n", id)
		},
	}

	tenantCmd.AddCommand(createCmd, listCmd, getCmd, deleteCmd)
	rootCmd.AddCommand(tenantCmd)
	return rootCmd
}
