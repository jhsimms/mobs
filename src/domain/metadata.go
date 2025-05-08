package domain

import (
	"fmt"
	"strings"
	"time"
)

// TenantMetadata extends Tenant with storage-specific properties
type TenantMetadata struct {
	// Core tenant properties
	TenantID  string       `json:"tenant_id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	Status    TenantStatus `json:"status"`

	// Storage-specific properties
	BucketName           string            `json:"bucket_name"`
	Region               string            `json:"region"`
	LastUpdatedAt        time.Time         `json:"last_updated_at"`
	Version              int64             `json:"version"`
	ProvisioningMetadata map[string]string `json:"provisioning_metadata"`
}

// NewTenantMetadata creates a new tenant metadata from a tenant
func NewTenantMetadata(tenant *Tenant, region string) (*TenantMetadata, error) {
	if tenant == nil {
		return nil, NewDomainError(ErrInvalidInput, "tenant cannot be nil", nil)
	}

	now := time.Now().UTC()
	bucketName := GenerateBucketName(tenant.TenantID, tenant.Name)

	metadata := &TenantMetadata{
		// Copy tenant properties
		TenantID:  tenant.TenantID,
		Name:      tenant.Name,
		CreatedAt: tenant.CreatedAt,
		Status:    tenant.Status,

		// Set storage-specific properties
		BucketName:           bucketName,
		Region:               region,
		LastUpdatedAt:        now,
		Version:              1,
		ProvisioningMetadata: make(map[string]string),
	}

	// Validate the new metadata
	if err := metadata.Validate(); err != nil {
		return nil, err
	}

	return metadata, nil
}

// ToTenant converts tenant metadata back to a tenant
func (tm *TenantMetadata) ToTenant() *Tenant {
	return &Tenant{
		TenantID:  tm.TenantID,
		Name:      tm.Name,
		CreatedAt: tm.CreatedAt,
		Status:    tm.Status,
	}
}

// Validate checks if the tenant metadata meets all validation rules
func (tm *TenantMetadata) Validate() error {
	var errors ValidationErrors

	// Validate tenant properties
	ValidateUUID("tenant_id", tm.TenantID, &errors)
	ValidateTenantName("name", tm.Name, &errors)
	ValidateTimestamp("created_at", tm.CreatedAt, &errors)
	ValidateTenantStatus("status", tm.Status, &errors)

	// Validate storage-specific properties
	ValidateBucketName("bucket_name", tm.BucketName, &errors)
	ValidateAWSRegion("region", tm.Region, &errors)
	ValidateTimestamp("last_updated_at", tm.LastUpdatedAt, &errors)
	ValidateVersion("version", tm.Version, &errors)

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// ChangeStatus transitions the tenant metadata to a new status
func (tm *TenantMetadata) ChangeStatus(newStatus TenantStatus) error {
	// Skip if status is not changing
	if tm.Status == newStatus {
		return nil
	}

	// Validate the status transition
	if !IsValidTransition(tm.Status, newStatus) {
		return NewDomainError(
			ErrInvalidTransition,
			fmt.Sprintf("Cannot transition from %s to %s", tm.Status, newStatus),
			nil,
		)
	}

	tm.Status = newStatus
	tm.LastUpdatedAt = time.Now().UTC()
	tm.Version++

	return nil
}

// IncrementVersion updates the version and last updated timestamp
func (tm *TenantMetadata) IncrementVersion() {
	tm.LastUpdatedAt = time.Now().UTC()
	tm.Version++
}

// SetProvisioningMetadata adds or updates a provisioning metadata entry
func (tm *TenantMetadata) SetProvisioningMetadata(key, value string) {
	if tm.ProvisioningMetadata == nil {
		tm.ProvisioningMetadata = make(map[string]string)
	}

	tm.ProvisioningMetadata[key] = value
	tm.IncrementVersion()
}

// GetProvisioningMetadata retrieves a provisioning metadata entry
func (tm *TenantMetadata) GetProvisioningMetadata(key string) (string, bool) {
	if tm.ProvisioningMetadata == nil {
		return "", false
	}

	value, exists := tm.ProvisioningMetadata[key]
	return value, exists
}

// String returns a string representation of the tenant metadata for logging
func (tm *TenantMetadata) String() string {
	return fmt.Sprintf(
		"TenantMetadata{ID: %s, Name: %s, Bucket: %s, Region: %s, Status: %s, Version: %d}",
		tm.TenantID, tm.Name, tm.BucketName, tm.Region, tm.Status, tm.Version,
	)
}

// GenerateBucketName creates a valid S3 bucket name from tenant properties
// Format: apt-{tenant-id-prefix}-{sanitized-name}
func GenerateBucketName(tenantID, name string) string {
	// Use first 8 chars of UUID
	idPrefix := ""
	if len(tenantID) >= 8 {
		idPrefix = tenantID[0:8]
	}

	// Sanitize name: lowercase, replace non-alphanumeric with hyphens
	sanitizedName := strings.ToLower(name)
	re := strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
	)
	sanitizedName = re.Replace(sanitizedName)

	// Trim sanitized name to fit bucket length limits (considering prefix length)
	maxNameLength := 63 - 13 // 63 is S3 limit, 13 is "apt-{8 chars}-" prefix (4+8+1=13)
	if len(sanitizedName) > maxNameLength {
		sanitizedName = sanitizedName[0:maxNameLength]
	}

	// Ensure we end with an alphanumeric character
	sanitizedName = strings.TrimRight(sanitizedName, "-")

	// Combine parts
	return fmt.Sprintf("apt-%s-%s", idPrefix, sanitizedName)
}
