package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a customer in the multi-tenant object storage system
type Tenant struct {
	// TODO: Evaluate properties
	TenantID  string       `json:"tenant_id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	Status    TenantStatus `json:"status"`
}

// NewTenant creates a new tenant with default values
func NewTenant(name string) (*Tenant, error) {
	tenantID := uuid.New().String()

	tenant := &Tenant{
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Status:    StatusProvisioning,
	}

	// Validate the new tenant
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	return tenant, nil
}

// Validate checks if the tenant meets all validation rules
func (t *Tenant) Validate() error {
	var errors ValidationErrors

	ValidateUUID("tenant_id", t.TenantID, &errors)
	ValidateTenantName("name", t.Name, &errors)
	ValidateTimestamp("created_at", t.CreatedAt, &errors)
	ValidateTenantStatus("status", t.Status, &errors)

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// ChangeStatus transitions the tenant to a new status
func (t *Tenant) ChangeStatus(newStatus TenantStatus) error {
	// Skip if status is not changing
	if t.Status == newStatus {
		return nil
	}

	// Validate the status transition
	if !IsValidTransition(t.Status, newStatus) {
		return NewDomainError(
			ErrInvalidTransition,
			fmt.Sprintf("Cannot transition from %s to %s", t.Status, newStatus),
			nil,
		)
	}

	t.Status = newStatus
	return nil
}

// String returns a string representation of the tenant for logging
func (t *Tenant) String() string {
	return fmt.Sprintf("Tenant{ID: %s, Name: %s, Status: %s}", t.TenantID, t.Name, t.Status)
}
