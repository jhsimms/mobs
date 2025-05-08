package cloverdb

import (
	"github.com/aptible/mobs/src/domain"
	"github.com/aptible/mobs/src/storage"
	"github.com/ostafen/clover"
)

const tenantsCollection = "tenants"

// CloverStore implements storage.TenantStore using CloverDB.
type CloverStore struct {
	db     *clover.DB
	dbPath string
}

// NewStore initializes a new CloverStore with the given dbPath.
func NewStore(dbPath string) (storage.TenantStore, error) {
	db, err := clover.Open(dbPath)
	if err != nil {
		return nil, err // TODO: Add better error handling and recovery
	}
	// Ensure the collection exists
	if ok, _ := db.HasCollection(tenantsCollection); !ok {
		db.CreateCollection(tenantsCollection)
	}
	return &CloverStore{db: db, dbPath: dbPath}, nil
}

// Create stores a new tenant and returns its metadata.
func (s *CloverStore) Create(tenant domain.Tenant) (*domain.TenantMetadata, error) {
	doc := clover.NewDocument()
	doc.Set("id", tenant.ID)
	doc.Set("name", tenant.Name)
	_, err := s.db.InsertOne(tenantsCollection, doc)
	if err != nil {
		return nil, err // TODO: Add proper error handling
	}
	return &domain.TenantMetadata{ID: tenant.ID, Name: tenant.Name}, nil
}

// Get retrieves tenant metadata by ID.
func (s *CloverStore) Get(tenantID string) (*domain.TenantMetadata, error) {
	filter := clover.Field("id").Eq(tenantID)
	docs, err := s.db.Query(tenantsCollection).Where(filter).FindAll()
	if err != nil || len(docs) == 0 {
		return nil, err // TODO: Return not found error
	}
	doc := docs[0]
	return &domain.TenantMetadata{
		ID:   doc.Get("id").(string),
		Name: doc.Get("name").(string),
	}, nil
}

// List returns all tenant metadata.
func (s *CloverStore) List() ([]domain.TenantMetadata, error) {
	docs, err := s.db.Query(tenantsCollection).FindAll()
	if err != nil {
		return nil, err // TODO: Add proper error handling
	}
	var tenants []domain.TenantMetadata
	for _, doc := range docs {
		id, _ := doc.Get("id").(string)
		name, _ := doc.Get("name").(string)
		tenants = append(tenants, domain.TenantMetadata{ID: id, Name: name})
	}
	return tenants, nil
}

// Delete removes a tenant by ID.
func (s *CloverStore) Delete(tenantID string) error {
	filter := clover.Field("id").Eq(tenantID)
	err := s.db.Query(tenantsCollection).Where(filter).Delete()
	if err != nil {
		return err // TODO: Add proper error handling
	}
	return nil
}
