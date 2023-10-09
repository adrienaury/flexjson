package json

type (
	Maker[T any] func() T
	Adder[T any] func(arr T, value any) T
	Keyer[T any] func(obj T, key string, value any) T
)

type ObjectStrategy[T any] interface {
	Make() T
	Add(obj T, key string, value any) T
}

type ArrayStrategy[T any] interface {
	Make() T
	Add(obj T, value any) T
}
