package flexjson_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/adrienaury/flexjson"
)

func TestDecodeStream(t *testing.T) {
	t.Parallel()

	const stream = `
		{"Surname": "Doe", "Name": "John", "Friends": [{"Surname": "Wizz", "Name": "Pat"}]}
`

	dec := flexjson.NewDecoder(strings.NewReader(stream))

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
