package json

type Object interface {
	Add(key string, value any)
	Get(key string) (any, bool)
	Has(key string) bool
	Len() int
	Keys() []string
	Value() []any
}

type Array interface {
	Append(value any)
	Get(index int) any
	Len() int
}

type (
	ObjectMaker func() Object
	ArrayMaker  func() Array
)

type TokenReader interface {
	More() bool
	Token() (Token, error)
}
