package v2

import (
	"encoding/json"

	"github.com/7sDream/geko"
)

type Object interface {
	Add(key string, value any)
}

type Array interface {
	Append(value ...any)
}

type (
	ObjectMaker func() Object
	ArrayMaker  func() Array
)

type (
	Token = json.Token
	Delim = json.Delim
)

type TokenReader interface {
	More() bool
	Token() (Token, error)
}

type Decoder struct {
	reader   TokenReader
	objmaker ObjectMaker
	arrmaker ArrayMaker
}

func NewDecoder(reader TokenReader) *Decoder {
	return &Decoder{
		reader: reader,
		objmaker: func() Object {
			return geko.NewMap[string, any]()
		},
		arrmaker: func() Array {
			return geko.NewList[any]()
		},
	}
}

func (d *Decoder) Decode() any {
	token, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	delim, isDelim := token.(Delim)

	switch {
	case isDelim && delim == Delim('{'):
		result := d.DecodeObject()
		d.assertNextToken('}')

		return result
	case isDelim && delim == Delim('['):
		result := d.DecodeArray()
		d.assertNextToken(']')

		return result
	case isDelim:
		panic("unexpected token" + string(delim))
	default:
		return token
	}
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

	value, err := d.reader.Token()
	if err != nil {
		panic(err)
	}

	delim, isDelim := value.(Delim)

	switch {
	case isDelim && delim == Delim('{'):
		obj.Add(keystr, d.DecodeObject())
		d.assertNextToken('}')
	case isDelim && delim == Delim('['):
		obj.Add(keystr, d.DecodeArray())
		d.assertNextToken(']')
	case isDelim:
		panic("unexpected token " + string(delim))
	default:
		obj.Add(keystr, value)
	}
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
