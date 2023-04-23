package any

import (
	"context"
	"golang.org/x/exp/constraints"
)

func GenericFunction[K []map[T]X, T int | bool, X ~string](x []K) T {
	var value T
	return value
}

type Repository[T any, ID any | string | constraints.Ordered | int | float32] interface {
	Save(entity T) T
}

type Controller[C context.Context, T any | int, Y ~int] struct {
	AnyField1 string
	AnyField2 int
}

func (c Controller[K, C, Y]) Index(ctx K, h C) {

}

type TestController struct {
	BaseController[int]
	Controller[context.Context, int16, int]
}

type Number interface {
	constraints.Ordered
	ToString()
}

type HttpHandler[C context.Context, K string | int, V constraints.Ordered | constraints.Complex, M ~string] func(ctx C, value V) K

type EventPublisher[E any, ID ~int] interface {
	Publish(e E)
}

type BaseController[M any] struct {
}
