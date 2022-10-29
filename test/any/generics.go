package any

import (
	"context"
	"golang.org/x/exp/constraints"
)

func GenericFunction[K []map[T]X, T int | bool, X ~string](x []K) T {
	var value T
	return value
}

type Repository[T, ID any] interface {
	Save(entity T) T
}

type Controller[C context.Context, T any] struct {
	AnyField1 string
	AnyField2 int
}

func (c Controller[K, C]) Index(ctx K, h C) {

}

type TestController struct {
	Controller[context.Context, int16]
}

type Number interface {
	constraints.Ordered
	ToString()
}

type HttpHandler[C context.Context, K string | int] func(ctx C) K

type EventPublisher[E any] interface {
	Publish(e E)
}
