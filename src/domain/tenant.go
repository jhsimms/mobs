package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system.
type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TenantMetadata contains metadata about a tenant.
type TenantMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// TODO: Add more metadata fields as needed
}

// NewTenant creates a new tenant with a generated ID
func NewTenant(name string) *Tenant {
	return &Tenant{
		ID:   uuid.New().String(),
		Name: name,
	}
}

// String returns a string representation of the tenant
func (t *Tenant) String() string {
	return fmt.Sprintf("Tenant{ID: %s, Name: %s}", t.ID, t.Name)
}
