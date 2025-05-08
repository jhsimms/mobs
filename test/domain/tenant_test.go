package domain_test

import (
	"testing"

	"github.com/aptible/mobs/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewTenant_ValidInput_CreatesSuccessfully(t *testing.T) {
	tenant := domain.NewTenant("test-tenant")
	assert.NotNil(t, tenant)
	assert.Equal(t, "test-tenant", tenant.Name)
	assert.NotEmpty(t, tenant.ID)
}
