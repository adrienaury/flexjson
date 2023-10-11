package flexjson

type (
	Maker[T any] func() T
	Adder[T any] func(arr T, value any) T
	Keyer[T any] func(obj T, key string, value any) T
)

type TokenReader interface {
	More() bool
	Token() (Token, error)
}
