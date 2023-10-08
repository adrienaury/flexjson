package v1

import (
	"fmt"

	"github.com/7sDream/geko"
)

func Decode(reader TokenReader) (Object, error) {
	token, err := reader.Token()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if delim, ok := token.(Delim); !ok || ok && delim != Delim('{') {
		return nil, fmt.Errorf("%w", ErrInvalidCharacter)
	}

	result := geko.NewMap[string, any]()
	for reader.More() {
		decode(reader, result)
	}

	token, err = reader.Token()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if delim, ok := token.(Delim); !ok || ok && delim != Delim('}') {
		return nil, fmt.Errorf("%w", ErrInvalidCharacter)
	}

	return result, nil
}

func decode(reader TokenReader, object Object) {
	key, err := reader.Token()
	if err != nil {
		panic(err)
	}

	value, err := reader.Token()
	if err != nil {
		panic(err)
	}

	if delim, ok := value.(Delim); ok && delim == Delim('{') {
		decoded := geko.NewMap[string, any]()
		for reader.More() {
			decode(reader, decoded)
		}
		object.Add(key.(string), decoded)
		reader.Token()
	} else if ok && delim == Delim('[') {
		decoded := []any{}
		for reader.More() {
			item, _ := Decode(reader)
			decoded = append(decoded, item)
		}
		object.Add(key.(string), decoded)
		reader.Token()
	} else if ok {
		panic(delim)
	} else {
		object.Add(key.(string), value)
	}
}
