package domain_test

import (
	"testing"
	"time"

	"github.com/aptible/mobs/src/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTenant_ValidInput_CreatesSuccessfully(t *testing.T) {
	// Test valid creation
	tenant, err := domain.NewTenant("test-tenant")

	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, "test-tenant", tenant.Name)
	assert.Equal(t, domain.StatusProvisioning, tenant.Status)
}

func TestNewTenant_ValidInput_GeneratesValidUUID(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	_, err = uuid.Parse(tenant.TenantID)
	assert.NoError(t, err)
}

func TestNewTenant_ValidInput_SetsCurrentTimestamp(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	diff := time.Since(tenant.CreatedAt)
	assert.Less(t, diff, 5*time.Second)
}

func TestNewTenant_NameTooShort_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("ab")

	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Contains(t, err.Error(), "must be between 3 and 64 characters")
}

func TestNewTenant_NameTooLong_ReturnsError(t *testing.T) {
	longName := "a123456789012345678901234567890123456789012345678901234567890123456789"
	tenant, err := domain.NewTenant(longName)

	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Contains(t, err.Error(), "must be between 3 and 64 characters")
}

func TestNewTenant_InvalidNameCharacters_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("invalid@name!")

	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Contains(t, err.Error(), "must contain only alphanumeric characters")
}

func TestTenant_Validate_ValidTenant_NoError(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	err = tenant.Validate()
	assert.NoError(t, err)
}

func TestTenant_Validate_InvalidUUID_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	tenant.TenantID = "not-a-uuid"
	err = tenant.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant_id: must be a valid UUID")
}

func TestTenant_Validate_EmptyUUID_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	tenant.TenantID = ""
	err = tenant.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant_id: cannot be empty")
}

func TestTenant_Validate_FutureCreatedAt_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	tenant.CreatedAt = time.Now().Add(24 * time.Hour)
	err = tenant.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "created_at: cannot be in the future")
}

func TestTenant_Validate_InvalidStatus_ReturnsError(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	tenant.Status = "INVALID_STATUS"
	err = tenant.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status: must be one of")
}

func TestTenant_ChangeStatus_ProvisioningToActive_Succeeds(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusProvisioning, tenant.Status)

	err = tenant.ChangeStatus(domain.StatusActive)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusActive, tenant.Status)
}

func TestTenant_ChangeStatus_ActiveToSuspended_Succeeds(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	// Setup: Change to active first
	err = tenant.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)

	// Test the transition to suspended
	err = tenant.ChangeStatus(domain.StatusSuspended)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusSuspended, tenant.Status)
}

func TestTenant_ChangeStatus_SuspendedToActive_Succeeds(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	// Setup: Change to active then suspended
	err = tenant.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)
	err = tenant.ChangeStatus(domain.StatusSuspended)
	assert.NoError(t, err)

	// Test the transition back to active
	err = tenant.ChangeStatus(domain.StatusActive)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusActive, tenant.Status)
}

func TestTenant_ChangeStatus_ActiveToProvisioning_Fails(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	// Setup: Change to active first
	err = tenant.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)

	// Test the invalid transition
	err = tenant.ChangeStatus(domain.StatusProvisioning)

	assert.Error(t, err)
	assert.Equal(t, domain.StatusActive, tenant.Status) // Should not change
	assert.Contains(t, err.Error(), "Cannot transition from ACTIVE to PROVISIONING")
}

func TestTenant_ChangeStatus_SameStatus_NoOp(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	// Setup: Change to active first
	err = tenant.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)

	// Test no-op transition (same status)
	err = tenant.ChangeStatus(domain.StatusActive)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusActive, tenant.Status)
}

func TestTenant_String_ReturnsFormattedString(t *testing.T) {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)

	str := tenant.String()

	assert.Contains(t, str, "Tenant{ID: ")
	assert.Contains(t, str, ", Name: test-tenant")
	assert.Contains(t, str, ", Status: PROVISIONING")
}
