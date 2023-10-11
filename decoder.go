package flexjson

import (
	"encoding/json"
	"io"
)

type (
	Token = json.Token
	Delim = json.Delim
)

type Decoder[O any, A any] struct {
	reader   TokenReader
	objMaker Maker[O]
	objKeyer Keyer[O]
	arrMaker Maker[A]
	arrAdder Adder[A]
}

func NewDecoderStandard(reader io.Reader) *Decoder[Object, Array] {
	return &Decoder[Object, Array]{
		reader:   json.NewDecoder(reader),
		objMaker: StandardObjectMaker(),
		objKeyer: StandardObjectAdder(),
		arrMaker: StandardArrayMaker(),
		arrAdder: StandardArrayAdder(),
	}
}

func NewDecoder[O any, A any](r TokenReader, om Maker[O], ok Keyer[O], am Maker[A], aadd Adder[A]) *Decoder[O, A] {
	return &Decoder[O, A]{
		reader:   r,
		objMaker: om,
		objKeyer: ok,
		arrMaker: am,
		arrAdder: aadd,
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
	object := d.objMaker()
	for d.reader.More() {
		object = d.decodeKeyValue(object)
	}

	return object
}

func (d *Decoder[O, A]) DecodeArray() A { //nolint:ireturn
	array := d.arrMaker()
	for d.reader.More() {
		array = d.arrAdder(array, d.Decode())
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

	return d.objKeyer(obj, keystr, d.Decode())
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
