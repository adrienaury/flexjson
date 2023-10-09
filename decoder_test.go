package json_test

import (
	"fmt"
	"strings"
	"testing"

	json "github.com/adrienaury/ordered-json"
)

func Test_DecodeStream(t *testing.T) {
	t.Parallel()

	const stream = `
		{"Surname": "Doe", "Name": "John"}
`

	dec := json.NewDecoder(strings.NewReader(stream))

	result := dec.Decode()

	object, ok := result.(json.Object)
	if !ok {
		t.Fatalf("result is not an object")
	}

	if len(object) != 2 {
		t.Fatalf("object len is invalid: %d", len(object))
	}

	fmt.Println(object)
}
