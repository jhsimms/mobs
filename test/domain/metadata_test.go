package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aptible/mobs/src/domain"
	"github.com/stretchr/testify/assert"
)

func createTestTenant(t *testing.T) *domain.Tenant {
	tenant, err := domain.NewTenant("test-tenant")
	assert.NoError(t, err)
	return tenant
}

func createTestMetadata(t *testing.T) *domain.TenantMetadata {
	tenant := createTestTenant(t)
	metadata, err := domain.NewTenantMetadata(tenant, "us-west-2")
	assert.NoError(t, err)
	return metadata
}

func TestNewTenantMetadata_ValidInput_CreatesSuccessfully(t *testing.T) {
	tenant := createTestTenant(t)
	region := "us-west-2"

	metadata, err := domain.NewTenantMetadata(tenant, region)

	assert.NoError(t, err)
	assert.NotNil(t, metadata)
}

func TestNewTenantMetadata_ValidInput_CopiesTenantProperties(t *testing.T) {
	tenant := createTestTenant(t)
	region := "us-west-2"

	metadata, err := domain.NewTenantMetadata(tenant, region)
	assert.NoError(t, err)

	assert.Equal(t, tenant.TenantID, metadata.TenantID)
	assert.Equal(t, tenant.Name, metadata.Name)
	assert.Equal(t, tenant.CreatedAt, metadata.CreatedAt)
	assert.Equal(t, tenant.Status, metadata.Status)
}

func TestNewTenantMetadata_ValidInput_SetsStorageProperties(t *testing.T) {
	tenant := createTestTenant(t)
	region := "us-west-2"

	metadata, err := domain.NewTenantMetadata(tenant, region)
	assert.NoError(t, err)

	assert.NotEmpty(t, metadata.BucketName)
	assert.Equal(t, region, metadata.Region)
	assert.NotZero(t, metadata.LastUpdatedAt)
	assert.Equal(t, int64(1), metadata.Version)
	assert.NotNil(t, metadata.ProvisioningMetadata)
}

func TestNewTenantMetadata_ValidInput_CreatesBucketWithCorrectFormat(t *testing.T) {
	tenant := createTestTenant(t)

	metadata, err := domain.NewTenantMetadata(tenant, "us-west-2")
	assert.NoError(t, err)

	assert.Contains(t, metadata.BucketName, "apt-")
	assert.Contains(t, metadata.BucketName, tenant.TenantID[0:8])
	assert.Contains(t, metadata.BucketName, "-test-tenant")
}

func TestNewTenantMetadata_NilTenant_ReturnsError(t *testing.T) {
	metadata, err := domain.NewTenantMetadata(nil, "us-west-2")

	assert.Error(t, err)
	assert.Nil(t, metadata)
	assert.Contains(t, err.Error(), "tenant cannot be nil")
}

func TestNewTenantMetadata_InvalidRegion_ReturnsError(t *testing.T) {
	tenant := createTestTenant(t)

	metadata, err := domain.NewTenantMetadata(tenant, "invalid-region")

	assert.Error(t, err)
	assert.Nil(t, metadata)
	assert.Contains(t, err.Error(), "region: must be a valid AWS region code")
}

func TestTenantMetadata_ToTenant_ConvertsCorrectly(t *testing.T) {
	original := createTestTenant(t)
	metadata, err := domain.NewTenantMetadata(original, "us-west-2")
	assert.NoError(t, err)

	converted := metadata.ToTenant()

	assert.Equal(t, original.TenantID, converted.TenantID)
	assert.Equal(t, original.Name, converted.Name)
	assert.Equal(t, original.CreatedAt, converted.CreatedAt)
	assert.Equal(t, original.Status, converted.Status)
}

func TestTenantMetadata_Validate_ValidMetadata_NoError(t *testing.T) {
	metadata := createTestMetadata(t)

	err := metadata.Validate()

	assert.NoError(t, err)
}

func TestTenantMetadata_Validate_InvalidBucketName_ReturnsError(t *testing.T) {
	metadata := createTestMetadata(t)
	metadata.BucketName = "invalid_bucket_name!"

	err := metadata.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket_name: must contain only lowercase alphanumeric characters")
}

func TestTenantMetadata_Validate_InvalidRegion_ReturnsError(t *testing.T) {
	metadata := createTestMetadata(t)
	metadata.Region = "invalid-region"

	err := metadata.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "region: must be a valid AWS region code")
}

func TestTenantMetadata_Validate_FutureLastUpdatedAt_ReturnsError(t *testing.T) {
	metadata := createTestMetadata(t)
	metadata.LastUpdatedAt = time.Now().Add(24 * time.Hour)

	err := metadata.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "last_updated_at: cannot be in the future")
}

func TestTenantMetadata_Validate_NegativeVersion_ReturnsError(t *testing.T) {
	metadata := createTestMetadata(t)
	metadata.Version = -1

	err := metadata.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version: must be non-negative")
}

func TestTenantMetadata_ChangeStatus_ProvisioningToActive_Succeeds(t *testing.T) {
	metadata := createTestMetadata(t)
	assert.Equal(t, domain.StatusProvisioning, metadata.Status)
	initialVersion := metadata.Version
	initialUpdated := metadata.LastUpdatedAt

	// Allow time to pass so timestamps will differ
	time.Sleep(10 * time.Millisecond)

	err := metadata.ChangeStatus(domain.StatusActive)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusActive, metadata.Status)
	assert.Equal(t, initialVersion+1, metadata.Version)
	assert.True(t, metadata.LastUpdatedAt.After(initialUpdated))
}

func TestTenantMetadata_ChangeStatus_InvalidTransition_Fails(t *testing.T) {
	metadata := createTestMetadata(t)

	// Setup: Change to active first
	err := metadata.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)

	initialVersion := metadata.Version
	initialUpdated := metadata.LastUpdatedAt

	// Allow time to pass so timestamps would differ if updated
	time.Sleep(10 * time.Millisecond)

	// Test invalid transition
	err = metadata.ChangeStatus(domain.StatusProvisioning)

	assert.Error(t, err)
	assert.Equal(t, domain.StatusActive, metadata.Status)   // Should not change
	assert.Equal(t, initialVersion, metadata.Version)       // Should not change
	assert.Equal(t, initialUpdated, metadata.LastUpdatedAt) // Should not change
}

func TestTenantMetadata_ChangeStatus_SameStatus_NoOp(t *testing.T) {
	metadata := createTestMetadata(t)

	// Setup: Change to active first
	err := metadata.ChangeStatus(domain.StatusActive)
	assert.NoError(t, err)

	initialVersion := metadata.Version

	// Test no-op transition (same status)
	err = metadata.ChangeStatus(domain.StatusActive)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusActive, metadata.Status)
	assert.Equal(t, initialVersion, metadata.Version) // Should not change
}

func TestTenantMetadata_IncrementVersion_UpdatesVersionAndTimestamp(t *testing.T) {
	metadata := createTestMetadata(t)
	initialVersion := metadata.Version
	initialUpdated := metadata.LastUpdatedAt

	// Allow time to pass so timestamps will differ
	time.Sleep(10 * time.Millisecond)

	metadata.IncrementVersion()

	assert.Equal(t, initialVersion+1, metadata.Version)
	assert.True(t, metadata.LastUpdatedAt.After(initialUpdated))
}

func TestTenantMetadata_SetProvisioningMetadata_AddsKeyAndIncrementsVersion(t *testing.T) {
	metadata := createTestMetadata(t)
	initialVersion := metadata.Version

	metadata.SetProvisioningMetadata("role_arn", "arn:aws:iam::123456789012:role/example-role")

	assert.Equal(t, initialVersion+1, metadata.Version)

	value, exists := metadata.GetProvisioningMetadata("role_arn")
	assert.True(t, exists)
	assert.Equal(t, "arn:aws:iam::123456789012:role/example-role", value)
}

func TestTenantMetadata_GetProvisioningMetadata_NonExistentKey_ReturnsNotExists(t *testing.T) {
	metadata := createTestMetadata(t)

	value, exists := metadata.GetProvisioningMetadata("non_existent_key")

	assert.False(t, exists)
	assert.Empty(t, value)
}

func TestGenerateBucketName_SimpleName_CreatesValidName(t *testing.T) {
	tenantID := "9e142839-29f5-4a0a-a5df-5cfc25b1168e"
	tenantName := "test-tenant"

	bucketName := domain.GenerateBucketName(tenantID, tenantName)

	assert.Contains(t, bucketName, "apt-")
	assert.Contains(t, bucketName, "9e142839")
	assert.Contains(t, bucketName, "-test-tenant")
	assert.LessOrEqual(t, len(bucketName), 63)
	assert.GreaterOrEqual(t, len(bucketName), 3)
}

func TestGenerateBucketName_WithSpaces_ReplacesWithHyphens(t *testing.T) {
	tenantID := "9e142839-29f5-4a0a-a5df-5cfc25b1168e"
	tenantName := "test tenant name"

	bucketName := domain.GenerateBucketName(tenantID, tenantName)

	assert.Contains(t, bucketName, "-test-tenant-name")
	assert.NotContains(t, bucketName, " ")
}

func TestGenerateBucketName_WithUppercase_ConvertedToLowercase(t *testing.T) {
	tenantID := "9e142839-29f5-4a0a-a5df-5cfc25b1168e"
	tenantName := "TestTenant"

	bucketName := domain.GenerateBucketName(tenantID, tenantName)

	assert.Contains(t, bucketName, "-testtenant")
	assert.NotContains(t, bucketName, "T")
}

func TestGenerateBucketName_WithUnderscores_ConvertedToHyphens(t *testing.T) {
	tenantID := "9e142839-29f5-4a0a-a5df-5cfc25b1168e"
	tenantName := "test_tenant_name"

	bucketName := domain.GenerateBucketName(tenantID, tenantName)

	assert.Contains(t, bucketName, "-test-tenant-name")
	assert.NotContains(t, bucketName, "_")
}

func TestGenerateBucketName_LongName_Truncates(t *testing.T) {
	tenantID := "9e142839-29f5-4a0a-a5df-5cfc25b1168e"
	tenantName := "this-is-a-very-long-tenant-name-that-should-get-truncated-because-it-exceeds-the-s3-bucket-name-length-limit"

	bucketName := domain.GenerateBucketName(tenantID, tenantName)

	// Debug info
	fmt.Printf("Bucket name: %s\n", bucketName)
	fmt.Printf("Bucket name length: %d\n", len(bucketName))

	assert.LessOrEqual(t, len(bucketName), 63)
	assert.Contains(t, bucketName, "apt-")
	assert.Contains(t, bucketName, "9e142839")
	// Use a smaller substring that would definitely be in the truncated name
	assert.Contains(t, bucketName, "-this-is-a-very-long-tenant-name")
}
