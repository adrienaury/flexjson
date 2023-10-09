package json

type Object = map[string]any

var DefaultObjectStrategy = StrategyObject[Object]{
	Make: func() Object { return make(Object) },
	Add:  func(obj Object, key string, value any) { obj[key] = value },
}

type Array = []any

var DefaultArrayStrategy = StrategyArray[Array]{
	Make: func() Array { return make(Array, 0) },
	Add:  func(arr Array, value any) { arr = append(arr, value) },
}
