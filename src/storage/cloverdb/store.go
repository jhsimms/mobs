package cloverdb

import (
	"fmt"

	"github.com/aptible/mobs/src/domain"
	"github.com/aptible/mobs/src/storage"
	"github.com/ostafen/clover"
)

// tenantsCollection is the name of the collection for tenants in CloverDB.
const tenantsCollection = "tenants"

// CloverStore implements storage.TenantStore using CloverDB as the backend.
type CloverStore struct {
	db     *clover.DB
	dbPath string
}

// NewStore initializes a new CloverStore with the given dbPath. Returns an error if the database cannot be opened.
func NewStore(dbPath string) (storage.TenantStore, error) {
	db, err := clover.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CloverDB at %s: %w", dbPath, err)
	}
	if ok, _ := db.HasCollection(tenantsCollection); !ok {
		if err := db.CreateCollection(tenantsCollection); err != nil {
			return nil, fmt.Errorf("failed to create collection %s: %w", tenantsCollection, err)
		}
	}
	return &CloverStore{db: db, dbPath: dbPath}, nil
}

// Create stores a new tenant and returns its metadata or an error.
func (s *CloverStore) Create(tenant domain.Tenant) (*domain.TenantMetadata, error) {
	doc := clover.NewDocument()
	doc.Set("id", tenant.ID)
	doc.Set("name", tenant.Name)
	_, err := s.db.InsertOne(tenantsCollection, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tenant: %w", err)
	}
	return &domain.TenantMetadata{ID: tenant.ID, Name: tenant.Name}, nil
}

// Get retrieves tenant metadata by ID or returns an error if not found.
func (s *CloverStore) Get(tenantID string) (*domain.TenantMetadata, error) {
	filter := clover.Field("id").Eq(tenantID)
	docs, err := s.db.Query(tenantsCollection).Where(filter).FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to query tenant by id: %w", err)
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("tenant not found: %s", tenantID)
	}
	doc := docs[0]
	id, _ := doc.Get("id").(string)
	name, _ := doc.Get("name").(string)
	return &domain.TenantMetadata{ID: id, Name: name}, nil
}

// List returns all tenant metadata or an error.
func (s *CloverStore) List() ([]domain.TenantMetadata, error) {
	docs, err := s.db.Query(tenantsCollection).FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	var tenants []domain.TenantMetadata
	for _, doc := range docs {
		id, _ := doc.Get("id").(string)
		name, _ := doc.Get("name").(string)
		tenants = append(tenants, domain.TenantMetadata{ID: id, Name: name})
	}
	return tenants, nil
}

// Delete removes a tenant by ID or returns an error if not found.
func (s *CloverStore) Delete(tenantID string) error {
	filter := clover.Field("id").Eq(tenantID)
	err := s.db.Query(tenantsCollection).Where(filter).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}
