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
	objStrategy ObjectStrategy[O]
	arrStrategy ArrayStrategy[A]
}

func NewDecoderStandard(reader io.Reader) *Decoder[Object, Array] {
	return &Decoder[Object, Array]{
		reader:      json.NewDecoder(reader),
		objStrategy: NewStandardObjectStrategy(),
		arrStrategy: NewStandardArrayStrategy(),
	}
}

func NewDecoder[O any, A any](reader TokenReader) *Decoder[O, A] {
	return &Decoder[O, A]{
		reader:      reader,
		objStrategy: nil,
		arrStrategy: nil,
	}
}

func (d *Decoder[O, A]) WithObjectStrategy(objStrategy ObjectStrategy[O]) {
	d.objStrategy = objStrategy
}

func (d *Decoder[O, A]) WithArrayStrategy(arrStrategy ArrayStrategy[A]) {
	d.arrStrategy = arrStrategy
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
		object = d.decodeKeyValue(object)
	}

	return object
}

func (d *Decoder[O, A]) DecodeArray() A { //nolint:ireturn
	array := d.arrStrategy.Make()
	for d.reader.More() {
		array = d.arrStrategy.Add(array, d.Decode())
	}

	return array
}

func (d *Decoder[O, A]) decodeKeyValue(obj O) O { //nolint:ireturn
	key, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	keystr, keyIsString := key.(string)
	if !keyIsString {
		panic("invalid key")
	}

	return d.objStrategy.Add(obj, keystr, d.Decode())
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
