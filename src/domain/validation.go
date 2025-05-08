package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Regular expressions for validation
var (
	// TODO: This is a temporary regex for tenant names. We should use a more robust solution.
	tenantNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{2,63}$`)

	// TODO: This is a temporary regex for bucket names. We should use a more robust solution.
	bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

	// AWS region code validation: standard AWS region format
	awsRegionRegex = regexp.MustCompile(`^[a-z]{2}-[a-z]+-\d$`)
)

const (
	errEmptyValue        = "cannot be empty"
	errInvalidUUID       = "must be a valid UUID"
	errTenantNameLength  = "must be between 3 and 64 characters"
	errTenantNameChars   = "must contain only alphanumeric characters, hyphens, and underscores"
	errTimestampFuture   = "cannot be in the future"
	errBucketNameLength  = "must be between 3 and 63 characters"
	errBucketNameChars   = "must contain only lowercase alphanumeric characters, periods, and hyphens"
	errBucketNamePrefix  = "cannot start with 'xn--'"
	errBucketNamePeriods = "cannot contain consecutive periods"
	errBucketNameAdj     = "cannot contain adjacent periods and hyphens"
	errAWSRegionFormat   = "must be a valid AWS region code"
	errVersionNegative   = "must be non-negative"
)

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, errEmptyValue)
		return
	}

	_, err := uuid.Parse(value)
	if err != nil {
		errors.Add(field, errInvalidUUID)
	}
}

// ValidateTenantName checks if a tenant name is valid
func ValidateTenantName(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, errEmptyValue)
		return
	}

	if len(value) < 3 || len(value) > 64 {
		errors.Add(field, errTenantNameLength)
		return
	}

	if !tenantNameRegex.MatchString(value) {
		errors.Add(field, errTenantNameChars)
	}
}

// ValidateTimestamp checks if a timestamp is valid (not in the future)
func ValidateTimestamp(field string, value time.Time, errors *ValidationErrors) {
	if value.IsZero() {
		errors.Add(field, errEmptyValue)
		return
	}

	if value.After(time.Now().UTC()) {
		errors.Add(field, errTimestampFuture)
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
		errors.Add(field, "must be one of [PROVISIONING ACTIVE SUSPENDED]")
	}
}

// ValidateBucketName checks if a bucket name follows S3 naming conventions
func ValidateBucketName(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, errEmptyValue)
		return
	}

	if len(value) < 3 || len(value) > 63 {
		errors.Add(field, errBucketNameLength)
		return
	}

	if strings.HasPrefix(value, "xn--") {
		errors.Add(field, errBucketNamePrefix)
		return
	}

	if strings.Contains(value, "..") {
		errors.Add(field, errBucketNamePeriods)
		return
	}

	if strings.Contains(value, ".-") || strings.Contains(value, "-.") {
		errors.Add(field, errBucketNameAdj)
		return
	}

	if !bucketNameRegex.MatchString(value) {
		errors.Add(field, errBucketNameChars)
	}
}

// ValidateAWSRegion checks if a region code is valid
func ValidateAWSRegion(field, value string, errors *ValidationErrors) {
	if value == "" {
		errors.Add(field, errEmptyValue)
		return
	}

	if !awsRegionRegex.MatchString(value) {
		errors.Add(field, errAWSRegionFormat)
	}
}

// ValidateVersion checks if a version number is valid
func ValidateVersion(field string, value int64, errors *ValidationErrors) {
	if value < 0 {
		errors.Add(field, errVersionNegative)
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
