package json

type StrategyObject[T any] struct {
	Make func() T
	Add  func(obj T, key string, value any)
}

type StrategyArray[T any] struct {
	Make func() T
	Add  func(arr T, value any)
}
