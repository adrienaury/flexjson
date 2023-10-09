package json

type Object = map[string]any

type objectStrategy[T any] struct {
	Maker Maker[T]
	Adder Keyer[T]
}

func (s objectStrategy[T]) Make() T { //nolint:ireturn
	return s.Maker()
}

func (s objectStrategy[T]) Add(obj T, key string, value any) T { //nolint:ireturn
	return s.Adder(obj, key, value)
}

func NewStandardObjectStrategy() ObjectStrategy[Object] {
	return objectStrategy[Object]{
		Maker: func() Object { return make(Object) },
		Adder: func(obj Object, key string, value any) Object {
			obj[key] = value

			return obj
		},
	}
}

type Array = []any

type arrayStrategy[T any] struct {
	Maker Maker[T]
	Adder Adder[T]
}

func (s arrayStrategy[T]) Make() T { //nolint:ireturn
	return s.Maker()
}

func (s arrayStrategy[T]) Add(obj T, value any) T { //nolint:ireturn
	return s.Adder(obj, value)
}

func NewStandardArrayStrategy() ArrayStrategy[Array] {
	return arrayStrategy[Array]{
		Maker: func() Array { return make(Array, 0) },
		Adder: func(arr Array, value any) Array { return append(arr, value) },
	}
}
