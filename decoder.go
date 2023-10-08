package json

import (
	"encoding/json"
	"io"
)

type (
	Token = json.Token
	Delim = json.Delim
)

type Decoder struct {
	reader   TokenReader
	objmaker ObjectMaker
	arrmaker ArrayMaker
}

func NewDecoderStandard(reader io.Reader) *Decoder {
	return &Decoder{
		reader:   json.NewDecoder(reader),
		objmaker: newObject,
		arrmaker: newArray,
	}
}

func NewDecoderStandardCustom(reader io.Reader, objmaker ObjectMaker, arrmaker ArrayMaker) *Decoder {
	return &Decoder{
		reader:   json.NewDecoder(reader),
		objmaker: objmaker,
		arrmaker: arrmaker,
	}
}

func NewDecoder(reader TokenReader) *Decoder {
	return &Decoder{
		reader:   reader,
		objmaker: newObject,
		arrmaker: newArray,
	}
}

func NewDecoderCustom(reader TokenReader, objmaker ObjectMaker, arrmaker ArrayMaker) *Decoder {
	return &Decoder{
		reader:   reader,
		objmaker: objmaker,
		arrmaker: arrmaker,
	}
}

func (d *Decoder) Decode() any {
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

func (d *Decoder) DecodeObject() Object { //nolint:ireturn
	object := d.objmaker()
	for d.reader.More() {
		d.decodeKeyValue(object)
	}

	return object
}

func (d *Decoder) DecodeArray() Array { //nolint:ireturn
	array := d.arrmaker()
	for d.reader.More() {
		array.Append(d.Decode())
	}

	return array
}

func (d *Decoder) decodeKeyValue(obj Object) {
	key, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	keystr, keyIsString := key.(string)
	if !keyIsString {
		panic("invalid key")
	}

	obj.Add(keystr, d.Decode())
}

func (d *Decoder) assertNextToken(is rune) {
	if token, err := d.reader.Token(); err != nil {
		panic(err)
	} else if delim, isDelim := token.(Delim); !isDelim {
		panic("unexpected token")
	} else if delim != Delim(is) {
		panic("unexpected token " + string(delim))
	}
}
