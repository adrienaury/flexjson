package v1_test

import (
	"strings"
	"testing"

	"github.com/7sDream/geko"
	v1 "github.com/adrienaury/ordered-json/pkg/v1"
	"github.com/goccy/go-json"
)

func BenchmarkV1(b *testing.B) {
	json := []byte(`{"a": 1, "b": 2, "a": 3}`)

	for n := 0; n < b.N; n++ {
		result, _ := geko.JSONUnmarshal(json, geko.UseObject(), geko.ObjectOnDuplicatedKey(geko.Ignore))
		_ = result.(geko.Object)
	}
}

func BenchmarkRef(b *testing.B) {
	data := []byte(`{"a": 1, "b": 2, "a": 3}`)

	for n := 0; n < b.N; n++ {
		// decoder := json.NewDecoder(strings.NewReader(string(data)))
		result := map[string]any{}
		// decoder.Decode(&result)
		json.Unmarshal(data, &result)
	}
}

func BenchmarkTest(b *testing.B) {
	data := []byte(`{"a": 1, "b": 2, "a": 3}`)

	for n := 0; n < b.N; n++ {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		result := geko.NewMap[string, any]()
		decoder.Decode(&result)
		// fmt.Println(result)
	}
}

func BenchmarkToken(b *testing.B) {
	data := []byte(`{"a": 1, "b": 2, "a": 3}`)

	for n := 0; n < b.N; n++ {
		decoder := json.NewDecoder(strings.NewReader(string(data)))

		result := geko.NewMap[string, any]()
		_, _ = decoder.Token()

		key, _ := decoder.Token()
		value, _ := decoder.Token()

		result.Add(key.(string), value)
	}
}

func BenchmarkTarget(b *testing.B) {
	data := []byte(`{"a": 1, "b": 2, "a": 3}`)

	for n := 0; n < b.N; n++ {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		result, _ := v1.Decode(decoder)
		result.Len()
		// fmt.Println(result)
	}
}
