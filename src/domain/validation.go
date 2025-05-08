package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Regular expressions for validation
var (
	// Tenant name validation: 3-64 chars, alphanumeric with hyphens and underscores
	tenantNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,63}$`)

	// S3 bucket name validation: 3-63 chars, lowercase alphanumeric and hyphens
	// Must start and end with lowercase letter or number
	// Cannot have consecutive hyphens or start with xn--
	bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

	// AWS region code validation: standard AWS region format
	awsRegionRegex = regexp.MustCompile(`^[a-z]{2}-[a-z]+-\d$`)
)

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, "cannot be empty")
		return
	}

	_, err := uuid.Parse(value)
	if err != nil {
		errors.Add(field, "must be a valid UUID")
	}
}

// ValidateTenantName checks if a tenant name is valid
func ValidateTenantName(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, "cannot be empty")
		return
	}

	if len(value) < 3 || len(value) > 64 {
		errors.Add(field, "must be between 3 and 64 characters")
		return
	}

	if !tenantNameRegex.MatchString(value) {
		errors.Add(field, "must contain only alphanumeric characters, hyphens, and underscores")
	}
}

// ValidateTimestamp checks if a timestamp is valid (not in the future)
func ValidateTimestamp(field string, value time.Time, errors *ValidationErrors) {
	if value.IsZero() {
		errors.Add(field, "cannot be empty")
		return
	}

	now := time.Now().UTC()
	if value.After(now) {
		errors.Add(field, "cannot be in the future")
	}
}

// ValidateTenantStatus checks if a tenant status is valid
func ValidateTenantStatus(field string, value TenantStatus, errors *ValidationErrors) {
	valid := false
	for _, status := range []TenantStatus{StatusProvisioning, StatusActive, StatusSuspended} {
		if value == status {
			valid = true
			break
		}
	}

	if !valid {
		errors.Add(field, fmt.Sprintf("must be one of %v", []TenantStatus{StatusProvisioning, StatusActive, StatusSuspended}))
	}
}

// ValidateBucketName checks if a bucket name follows S3 naming conventions
func ValidateBucketName(field, value string, errors *ValidationErrors) {
	// TODO: Improve this validation, look for something off-the-shelf
	if value == "" {
		errors.Add(field, "cannot be empty")
		return
	}

	if len(value) < 3 || len(value) > 63 {
		errors.Add(field, "must be between 3 and 63 characters")
		return
	}

	if strings.HasPrefix(value, "xn--") {
		errors.Add(field, "cannot start with 'xn--'")
		return
	}

	if strings.Contains(value, "..") {
		errors.Add(field, "cannot contain consecutive periods")
		return
	}

	if strings.Contains(value, ".-") || strings.Contains(value, "-.") {
		errors.Add(field, "cannot contain adjacent periods and hyphens")
		return
	}

	if !bucketNameRegex.MatchString(value) {
		errors.Add(field, "must contain only lowercase alphanumeric characters, periods, and hyphens")
	}
}

// ValidateAWSRegion checks if a region code is valid
func ValidateAWSRegion(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, "cannot be empty")
		return
	}

	if !awsRegionRegex.MatchString(value) {
		errors.Add(field, "must be a valid AWS region code")
	}
}

// ValidateVersion checks if a version number is valid
func ValidateVersion(field string, value int64, errors *ValidationErrors) {
	if value < 0 {
		errors.Add(field, "must be non-negative")
	}
}

// TenantStatus represents the current state of a tenant
type TenantStatus string

const (
	// StatusProvisioning indicates tenant resources are being set up
	StatusProvisioning TenantStatus = "PROVISIONING"

	// StatusActive indicates the tenant is fully provisioned and operational
	StatusActive TenantStatus = "ACTIVE"

	// StatusSuspended indicates the tenant is temporarily disabled
	StatusSuspended TenantStatus = "SUSPENDED"
)

// IsValidTransition checks if a status transition is allowed
func IsValidTransition(from, to TenantStatus) bool {
	// Define allowed transitions
	allowedTransitions := map[TenantStatus][]TenantStatus{
		StatusProvisioning: {StatusActive, StatusSuspended},
		StatusActive:       {StatusSuspended},
		StatusSuspended:    {StatusActive},
	}

	// Check if transition is allowed
	allowed, exists := allowedTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}
