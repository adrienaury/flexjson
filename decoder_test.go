package flexjson_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adrienaury/flexjson"
	goccy "github.com/goccy/go-json"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestDecodeStream(t *testing.T) {
	t.Parallel()

	const stream = `
		{"Surname": "Doe", "Name": "John", "Friends": [{"Surname": "Wizz", "Name": "Pat"}]}
`

	dec := flexjson.NewDecoderStandard(strings.NewReader(stream))

	result := dec.Decode()

	object, ok := result.(flexjson.Object)
	if !ok {
		t.Fatalf("result is not an object")
	}

	if len(object) != 3 {
		t.Fatalf("object len is invalid: %d", len(object))
	}

	fmt.Println(object)
}

func TestCustomStrategy(t *testing.T) {
	t.Parallel()

	const stream = `
		{"Surname": "Doe", "Name": "John", "Friends": [{"Surname": "Wizz", "Name": "Pat"}]}
`

	dec := flexjson.NewDecoder(
		goccy.NewDecoder(strings.NewReader(stream)),
		func() *orderedmap.OrderedMap[string, any] { return orderedmap.New[string, any]() },
		func(obj *orderedmap.OrderedMap[string, any], key string, value any) *orderedmap.OrderedMap[string, any] {
			obj.Set(key, value)

			return obj
		},
		flexjson.StandardArrayMaker(),
		flexjson.StandardArrayAdder(),
	)

	result := dec.Decode()

	fmt.Println(result)
}
