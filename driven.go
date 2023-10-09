package json

type TokenReader interface {
	More() bool
	Token() (Token, error)
}
