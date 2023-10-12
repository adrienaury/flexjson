package flexjson

type (
	Maker[T any] func() (T, error)
	Adder[T any] func(arr T, value any) (T, error)
	Keyer[T any] func(obj T, key string, value any) (T, error)
)

type TokenReader interface {
	More() bool
	Token() (Token, error)
}
