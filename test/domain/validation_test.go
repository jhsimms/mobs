package domain_test

import (
	"testing"
	"time"

	"github.com/aptible/mobs/src/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateUUID(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid UUID",
			value:       uuid.New().String(),
			expectError: false,
		},
		{
			name:          "Empty UUID",
			value:         "",
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:          "Invalid UUID",
			value:         "not-a-uuid",
			expectError:   true,
			errorContains: "must be a valid UUID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateUUID("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateTenantName(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid name",
			value:       "test-tenant",
			expectError: false,
		},
		{
			name:          "Empty name",
			value:         "",
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:          "Too short",
			value:         "ab",
			expectError:   true,
			errorContains: "must be between 3 and 64 characters",
		},
		{
			name:          "Too long",
			value:         "a123456789012345678901234567890123456789012345678901234567890123456789",
			expectError:   true,
			errorContains: "must be between 3 and 64 characters",
		},
		{
			name:          "Invalid characters",
			value:         "test@tenant",
			expectError:   true,
			errorContains: "must contain only alphanumeric characters",
		},
		{
			name:        "With underscores",
			value:       "test_tenant_name",
			expectError: false,
		},
		{
			name:        "With hyphens",
			value:       "test-tenant-name",
			expectError: false,
		},
		{
			name:        "With numbers",
			value:       "test-tenant-123",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateTenantName("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateTimestamp(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	testCases := []struct {
		name          string
		value         time.Time
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid timestamp (now)",
			value:       now,
			expectError: false,
		},
		{
			name:        "Valid timestamp (past)",
			value:       past,
			expectError: false,
		},
		{
			name:          "Zero timestamp",
			value:         time.Time{},
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:          "Future timestamp",
			value:         future,
			expectError:   true,
			errorContains: "cannot be in the future",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateTimestamp("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateTenantStatus(t *testing.T) {
	testCases := []struct {
		name          string
		value         domain.TenantStatus
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid status (PROVISIONING)",
			value:       domain.StatusProvisioning,
			expectError: false,
		},
		{
			name:        "Valid status (ACTIVE)",
			value:       domain.StatusActive,
			expectError: false,
		},
		{
			name:        "Valid status (SUSPENDED)",
			value:       domain.StatusSuspended,
			expectError: false,
		},
		{
			name:          "Invalid status",
			value:         "INVALID_STATUS",
			expectError:   true,
			errorContains: "must be one of",
		},
		{
			name:          "Empty status",
			value:         "",
			expectError:   true,
			errorContains: "must be one of",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateTenantStatus("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateBucketName(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid bucket name",
			value:       "apt-12345678-test-tenant",
			expectError: false,
		},
		{
			name:          "Empty bucket name",
			value:         "",
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:          "Too short",
			value:         "ab",
			expectError:   true,
			errorContains: "must be between 3 and 63 characters",
		},
		{
			name:          "Too long",
			value:         "a123456789012345678901234567890123456789012345678901234567890123456789",
			expectError:   true,
			errorContains: "must be between 3 and 63 characters",
		},
		{
			name:          "With uppercase",
			value:         "Apt-12345678-test",
			expectError:   true,
			errorContains: "must contain only lowercase alphanumeric characters",
		},
		{
			name:          "With underscores",
			value:         "apt_12345678_test",
			expectError:   true,
			errorContains: "must contain only lowercase alphanumeric characters",
		},
		{
			name:          "Starts with xn--",
			value:         "xn--test-bucket",
			expectError:   true,
			errorContains: "cannot start with 'xn--'",
		},
		{
			name:          "With consecutive periods",
			value:         "apt..test",
			expectError:   true,
			errorContains: "cannot contain consecutive periods",
		},
		{
			name:          "With period and hyphen adjacent",
			value:         "apt.-test",
			expectError:   true,
			errorContains: "cannot contain adjacent periods and hyphens",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateBucketName("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateAWSRegion(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid region (us-east-1)",
			value:       "us-east-1",
			expectError: false,
		},
		{
			name:        "Valid region (us-west-2)",
			value:       "us-west-2",
			expectError: false,
		},
		{
			name:        "Valid region (eu-central-1)",
			value:       "eu-central-1",
			expectError: false,
		},
		{
			name:          "Empty region",
			value:         "",
			expectError:   true,
			errorContains: "cannot be empty",
		},
		{
			name:          "Invalid format",
			value:         "useast1",
			expectError:   true,
			errorContains: "must be a valid AWS region code",
		},
		{
			name:          "Invalid format with spaces",
			value:         "us east 1",
			expectError:   true,
			errorContains: "must be a valid AWS region code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateAWSRegion("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	testCases := []struct {
		name          string
		value         int64
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid version (0)",
			value:       0,
			expectError: false,
		},
		{
			name:        "Valid version (positive)",
			value:       42,
			expectError: false,
		},
		{
			name:          "Invalid version (negative)",
			value:         -1,
			expectError:   true,
			errorContains: "must be non-negative",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errors domain.ValidationErrors
			domain.ValidateVersion("test_field", tc.value, &errors)

			if tc.expectError {
				assert.True(t, errors.HasErrors())
				assert.Contains(t, errors.Error(), tc.errorContains)
			} else {
				assert.False(t, errors.HasErrors())
			}
		})
	}
}

func TestValidationErrorsCollection(t *testing.T) {
	var errors domain.ValidationErrors

	// Should start with no errors
	assert.False(t, errors.HasErrors())
	assert.Empty(t, errors.Error())

	// Add a single error
	errors.Add("field1", "error message 1")
	assert.True(t, errors.HasErrors())
	assert.Contains(t, errors.Error(), "validation failed: field1: error message 1")

	// Add multiple errors
	errors.Add("field2", "error message 2")
	errors.Add("field3", "error message 3")

	// Check error string formatting
	errStr := errors.Error()
	assert.Contains(t, errStr, "field1: error message 1")
	assert.Contains(t, errStr, "field2: error message 2")
	assert.Contains(t, errStr, "field3: error message 3")
	assert.Contains(t, errStr, ";") // Delimiter between errors
}

func TestIsValidTransition(t *testing.T) {
	testCases := []struct {
		name    string
		from    domain.TenantStatus
		to      domain.TenantStatus
		isValid bool
	}{
		{
			name:    "PROVISIONING -> ACTIVE",
			from:    domain.StatusProvisioning,
			to:      domain.StatusActive,
			isValid: true,
		},
		{
			name:    "PROVISIONING -> SUSPENDED",
			from:    domain.StatusProvisioning,
			to:      domain.StatusSuspended,
			isValid: true,
		},
		{
			name:    "ACTIVE -> SUSPENDED",
			from:    domain.StatusActive,
			to:      domain.StatusSuspended,
			isValid: true,
		},
		{
			name:    "SUSPENDED -> ACTIVE",
			from:    domain.StatusSuspended,
			to:      domain.StatusActive,
			isValid: true,
		},
		{
			name:    "ACTIVE -> PROVISIONING",
			from:    domain.StatusActive,
			to:      domain.StatusProvisioning,
			isValid: false,
		},
		{
			name:    "SUSPENDED -> PROVISIONING",
			from:    domain.StatusSuspended,
			to:      domain.StatusProvisioning,
			isValid: false,
		},
		{
			name:    "Invalid status -> any",
			from:    "INVALID_STATUS",
			to:      domain.StatusActive,
			isValid: false,
		},
		{
			name:    "Same status (ACTIVE -> ACTIVE)",
			from:    domain.StatusActive,
			to:      domain.StatusActive,
			isValid: false, // IsValidTransition checks actual transitions, not no-ops
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.IsValidTransition(tc.from, tc.to)
			assert.Equal(t, tc.isValid, result)
		})
	}
}
