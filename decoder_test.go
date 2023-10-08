package json_test

import (
	"strings"
	"testing"

	json "github.com/adrienaury/ordered-json"
)

func Test_DecodeStream(t *testing.T) {
	t.Parallel()

	const stream = `
		{"Surname": "Doe", "Name": "John"}
`

	dec := json.NewDecoderStandard(strings.NewReader(stream))

	result := dec.Decode()

	object, ok := result.(json.Object)
	if !ok {
		t.Fatalf("result is not an object")
	}

	if object.Len() != 2 {
		t.Fatalf("object len is invalid: %d", object.Len())
	}

	if object.Keys()[0] != "Surname" {
		t.Fatalf("object first key is invalid: %s", object.Keys()[0])
	}

	if object.Keys()[1] != "Name" {
		t.Fatalf("object second key is invalid")
	}
}
