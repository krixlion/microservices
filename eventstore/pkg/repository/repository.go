package repository

import "context"

// Repository CRUD abstraction for DB
type Repository[T any] interface {
	Create(context.Context, T) error
	Get(context.Context, string) (T, error)
	Index(context.Context) ([]T, error)
}
