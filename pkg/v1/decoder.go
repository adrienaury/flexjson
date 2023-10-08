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

	if delim, ok := value.(Delim); ok && delim != Delim('{') {
		decoded := geko.NewMap[string, any]()
		for reader.More() {
			decode(reader, decoded)
		}
	} else if ok {
		panic(ErrInvalidCharacter)
	} else {
		object.Add(key.(string), value)
	}
}
