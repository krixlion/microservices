package repository

// Repository CRUD abstraction for DB
type Repository[T, ID any] interface {
	Create(T) error
	Get(ID) (T, error)
	Index() []T
}
