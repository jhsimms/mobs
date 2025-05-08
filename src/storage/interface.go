package storage

import "github.com/aptible/mobs/src/domain"

// TenantStore defines the interface for tenant persistence operations.
type TenantStore interface {
	// Create stores a new tenant and returns its metadata.
	Create(tenant domain.Tenant) (*domain.TenantMetadata, error)
	// Get retrieves tenant metadata by ID.
	Get(tenantID string) (*domain.TenantMetadata, error)
	// List returns all tenant metadata.
	List() ([]domain.TenantMetadata, error)
	// Delete removes a tenant by ID.
	Delete(tenantID string) error
}
