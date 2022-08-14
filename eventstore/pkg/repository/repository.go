package repository

// Repository CRUD abstraction for DB
type Repository[T any] interface {
	Create(T) error
	Get(string) (T, error)
	Index(string) ([]T, error) // takes search params as an argument
}
