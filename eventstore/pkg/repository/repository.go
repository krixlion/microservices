package repository

// Repository CRUD abstraction for DB
type Repository[T, ID any] interface {
	Index() ([]T, error)
	Find(ID) (T, error)
	Create(T) error
	Update(T) error
	Delete(ID) error
}
