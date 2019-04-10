package recipe

// StorageGateway represent a data storage service
type StorageGateway interface {
	GetRange(start, limit uint64) ([]*Recipe, error)
	GetByID(id string) (*Recipe, error)
	DeleteByID(id string) error
	Store(recipe *Recipe) error
	Update(recipe *Recipe) error
	Delete(recipe *Recipe) error
	Search(pattern string) ([]*Recipe, error)
}
