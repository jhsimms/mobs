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
	// Additional metadata fields can be added here as needed.
}

// NewTenant creates a new Tenant with a generated unique ID and the provided name.
func NewTenant(name string) *Tenant {
	return &Tenant{
		ID:   uuid.New().String(),
		Name: name,
	}
}

// String returns a human-readable representation of the Tenant.
func (t *Tenant) String() string {
	return fmt.Sprintf("Tenant{ID: %s, Name: %s}", t.ID, t.Name)
}
