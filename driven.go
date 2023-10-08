package json

type Object interface {
	Add(key string, value any)
}

type Array interface {
	Append(value any)
}

type (
	ObjectMaker func() Object
	ArrayMaker  func() Array
)

type TokenReader interface {
	More() bool
	Token() (Token, error)
}
