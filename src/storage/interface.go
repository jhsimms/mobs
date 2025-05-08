package storage

import "github.com/aptible/mobs/src/domain"

// TenantStore defines the interface for tenant persistence operations.
type TenantStore interface {
	// Create stores a new tenant and returns its metadata or an error.
	Create(tenant domain.Tenant) (*domain.TenantMetadata, error)
	// Get retrieves tenant metadata by ID or returns an error if not found.
	Get(tenantID string) (*domain.TenantMetadata, error)
	// List returns all tenant metadata or an error.
	List() ([]domain.TenantMetadata, error)
	// Delete removes a tenant by ID or returns an error if not found.
	Delete(tenantID string) error
}
