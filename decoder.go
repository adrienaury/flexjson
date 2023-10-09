package json

import (
	"encoding/json"
	"io"
)

type (
	Token = json.Token
	Delim = json.Delim
)

type Decoder[O any, A any] struct {
	reader      TokenReader
	objStrategy StrategyObject[O]
	arrStrategy StrategyArray[A]
}

func NewDecoder(reader io.Reader) *Decoder[Object, Array] {
	return &Decoder[Object, Array]{
		reader:      json.NewDecoder(reader),
		objStrategy: DefaultObjectStrategy,
		arrStrategy: DefaultArrayStrategy,
	}
}

func (d *Decoder[O, A]) Decode() any {
	token, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	delim, isDelim := token.(Delim)

	var result any

	switch {
	case isDelim && delim == Delim('{'):
		result = d.DecodeObject()
		d.assertNextToken('}')
	case isDelim && delim == Delim('['):
		result = d.DecodeArray()
		d.assertNextToken(']')
	case isDelim:
		panic("unexpected token " + string(delim))
	default:
		result = token
	}

	return result
}

func (d *Decoder[O, A]) DecodeObject() O { //nolint:ireturn
	object := d.objStrategy.Make()
	for d.reader.More() {
		d.decodeKeyValue(object)
	}

	return object
}

func (d *Decoder[O, A]) DecodeArray() A { //nolint:ireturn
	array := d.arrStrategy.Make()
	for d.reader.More() {
		d.arrStrategy.Add(array, d.Decode())
	}

	return array
}

func (d *Decoder[O, A]) decodeKeyValue(obj O) {
	key, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	keystr, keyIsString := key.(string)
	if !keyIsString {
		panic("invalid key")
	}

	d.objStrategy.Add(obj, keystr, d.Decode())
}

func (d *Decoder[O, A]) assertNextToken(is rune) {
	if token, err := d.reader.Token(); err != nil {
		panic(err)
	} else if delim, isDelim := token.(Delim); !isDelim {
		panic("unexpected token")
	} else if delim != Delim(is) {
		panic("unexpected token " + string(delim))
	}
}
