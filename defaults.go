package flexjson

type (
	Object = map[string]any
	Array  = []any
)

func StandardObjectMaker() Maker[Object] {
	return func() (Object, error) { return make(Object), nil }
}

func StandardObjectAdder() Keyer[Object] {
	return func(obj Object, key string, value any) (Object, error) {
		obj[key] = value

		return obj, nil
	}
}

func StandardArrayMaker() Maker[Array] {
	return func() (Array, error) { return make(Array, 0), nil }
}

func StandardArrayAdder() Adder[Array] {
	return func(arr Array, value any) (Array, error) { return append(arr, value), nil }
}
