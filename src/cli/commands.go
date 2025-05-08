package main

import (
	"fmt"
	"os"

	"github.com/aptible/mobs/src/service"
	"github.com/spf13/cobra"
)

// setupCommands sets up the root and tenant commands.
func setupCommands(tenantService *service.TenantService) *cobra.Command {
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
				fmt.Fprintf(os.Stderr, "Failed to create tenant: %v\n", err)
				os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Failed to list tenants: %v\n", err)
				os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Failed to get tenant: %v\n", err)
				os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Failed to delete tenant: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Tenant deleted: %s\n", id)
		},
	}

	tenantCmd.AddCommand(createCmd, listCmd, getCmd, deleteCmd)
	rootCmd.AddCommand(tenantCmd)
	return rootCmd
}
