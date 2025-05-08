package service

import (
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

// CreateTenant creates a new tenant with the given name. Returns metadata or error if validation fails.
func (s *TenantService) CreateTenant(name string) (*domain.TenantMetadata, error) {
	errs := domain.ValidationErrors{}
	domain.ValidateTenantName("name", name, &errs)
	if errs.HasErrors() {
		return nil, errs
	}
	tenant := domain.NewTenant(name)
	return s.store.Create(*tenant)
}

// GetTenant retrieves tenant metadata by ID. Returns error if ID is invalid.
func (s *TenantService) GetTenant(id string) (*domain.TenantMetadata, error) {
	errs := domain.ValidationErrors{}
	domain.ValidateUUID("id", id, &errs)
	if errs.HasErrors() {
		return nil, errs
	}
	return s.store.Get(id)
}

// ListTenants returns all tenant metadata.
func (s *TenantService) ListTenants() ([]domain.TenantMetadata, error) {
	return s.store.List()
}

// DeleteTenant removes a tenant by ID. Returns error if ID is invalid.
func (s *TenantService) DeleteTenant(id string) error {
	errs := domain.ValidationErrors{}
	domain.ValidateUUID("id", id, &errs)
	if errs.HasErrors() {
		return errs
	}
	return s.store.Delete(id)
}
