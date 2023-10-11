package flexjson

type (
	Object = map[string]any
	Array  = []any
)

func StandardObjectMaker() Maker[Object] {
	return func() Object { return make(Object) }
}

func StandardObjectAdder() Keyer[Object] {
	return func(obj Object, key string, value any) Object {
		obj[key] = value

		return obj
	}
}

func StandardArrayMaker() Maker[Array] {
	return func() Array { return make(Array, 0) }
}

func StandardArrayAdder() Adder[Array] {
	return func(arr Array, value any) Array { return append(arr, value) }
}
