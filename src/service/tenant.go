package service

import (
	"errors"

	"github.com/aptible/mobs/src/domain"
	"github.com/aptible/mobs/src/storage"
)

// TenantService provides business logic for tenant operations.
type TenantService struct {
	store storage.TenantStore
}

// NewTenantService creates a new TenantService.
func NewTenantService(store storage.TenantStore) *TenantService {
	return &TenantService{store: store}
}

// CreateTenant creates a new tenant with the given name.
func (s *TenantService) CreateTenant(name string) (*domain.TenantMetadata, error) {
	if name == "" {
		return nil, errors.New("tenant name cannot be empty") // TODO: Add better validation
	}
	tenant := domain.NewTenant(name)
	return s.store.Create(*tenant)
}

// GetTenant retrieves tenant metadata by ID.
func (s *TenantService) GetTenant(id string) (*domain.TenantMetadata, error) {
	if id == "" {
		return nil, errors.New("tenant ID cannot be empty") // TODO: Add better validation
	}
	return s.store.Get(id)
}

// ListTenants returns all tenant metadata.
func (s *TenantService) ListTenants() ([]domain.TenantMetadata, error) {
	return s.store.List()
}

// DeleteTenant removes a tenant by ID.
func (s *TenantService) DeleteTenant(id string) error {
	if id == "" {
		return errors.New("tenant ID cannot be empty") // TODO: Add better validation
	}
	return s.store.Delete(id)
}
