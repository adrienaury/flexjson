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

func (a *array) Get(index int) any {
	return a.slice[index]
}

func (a *array) Len() int {
	return len(a.slice)
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
	o.order = append(o.order, key)
}

func (o *object) Get(key string) (any, bool) {
	v, has := o.inner[key]

	return v, has
}

func (o *object) Has(key string) bool {
	_, has := o.inner[key]

	return has
}

func (o *object) Len() int {
	return len(o.inner)
}

func (o *object) Keys() []string {
	keys := make([]string, o.Len())
	copy(keys, o.order)

	return keys
}

func (o *object) Value() []any {
	values := make([]any, 0, o.Len())
	for _, key := range o.order {
		values = append(values, o.inner[key])
	}

	return values
}
