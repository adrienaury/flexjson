package json

type array struct {
	slice []any
}

func newArray() Array { //nolint:ireturn
	return &array{
		slice: []any{},
	}
}

func (a *array) Append(value any) {
	a.slice = append(a.slice, value)
}

type object struct {
	order []string
	inner map[string]any
}

func newObject() Object { //nolint:ireturn
	return &object{
		order: []string{},
		inner: map[string]any{},
	}
}

func (o *object) Add(key string, value any) {
	if _, has := o.inner[key]; has {
		return
	}

	o.inner[key] = value
}
