package flexjson

import (
	"encoding/json"
	"fmt"
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

func NewDecoder(reader io.Reader) *Decoder[Object, Array] {
	return &Decoder[Object, Array]{
		reader:   json.NewDecoder(reader),
		objMaker: StandardObjectMaker(),
		objKeyer: StandardObjectAdder(),
		arrMaker: StandardArrayMaker(),
		arrAdder: StandardArrayAdder(),
	}
}

func NewFlexDecoder[O any, A any](
	reader TokenReader,
	objMaker Maker[O],
	objKeyer Keyer[O],
	arrMaker Maker[A],
	arrAdder Adder[A],
) *Decoder[O, A] {
	return &Decoder[O, A]{
		reader:   reader,
		objMaker: objMaker,
		objKeyer: objKeyer,
		arrMaker: arrMaker,
		arrAdder: arrAdder,
	}
}

func (d *Decoder[O, A]) Decode() (any, error) {
	token, err := d.reader.Token()
	if err != nil {
		return nil, err
	}

	delim, isDelim := token.(Delim)

	var result any

	switch {
	case isDelim && delim == Delim('{'):
		result, err = d.decodeObject()
		if err != nil {
			return nil, err
		}
		if err := d.assertNextToken('}'); err != nil {
			return nil, err
		}
	case isDelim && delim == Delim('['):
		result, err = d.decodeArray()
		if err != nil {
			return nil, err
		}
		if err := d.assertNextToken(']'); err != nil {
			return nil, err
		}
	case isDelim:
		return nil, fmt.Errorf("unexpected token %s", string(delim))
	default:
		result = token
	}

	return result, nil
}

func (d *Decoder[O, A]) DecodeObject() (O, error) { //nolint:ireturn
	if err := d.assertNextToken('{'); err != nil {
		return *new(O), err
	}

	var result O
	var err error

	result, err = d.decodeObject()
	if err != nil {
		return *new(O), err
	}

	if err := d.assertNextToken('}'); err != nil {
		return *new(O), err
	}

	return result, nil
}

func (d *Decoder[O, A]) decodeObject() (O, error) { //nolint:ireturn
	object, err := d.objMaker()
	if err != nil {
		return *new(O), err
	}

	for d.reader.More() {
		if object, err = d.decodeKeyValue(object); err != nil {
			return *new(O), err
		}
	}

	return object, nil
}

func (d *Decoder[O, A]) decodeKeyValue(obj O) (O, error) { //nolint:ireturn
	key, err := d.reader.Token()
	if err != nil {
		return *new(O), err
	}

	keystr, keyIsString := key.(string)
	if !keyIsString {
		return *new(O), fmt.Errorf("invalid key")
	}

	result, err := d.Decode()
	if err != nil {
		return *new(O), err
	}

	return d.objKeyer(obj, keystr, result)
}

func (d *Decoder[O, A]) DecodeArray() (A, error) { //nolint:ireturn
	if err := d.assertNextToken('['); err != nil {
		return *new(A), err
	}

	var array A
	var err error

	array, err = d.decodeArray()
	if err != nil {
		return *new(A), err
	}

	if err := d.assertNextToken(']'); err != nil {
		return *new(A), err
	}

	return array, nil
}

func (d *Decoder[O, A]) decodeArray() (A, error) { //nolint:ireturn
	array, err := d.arrMaker()
	if err != nil {
		return *new(A), err
	}

	for d.reader.More() {
		item, err := d.Decode()
		if err != nil {
			return *new(A), err
		}

		array, err = d.arrAdder(array, item)
		if err != nil {
			return *new(A), err
		}
	}

	return array, nil
}

func (d *Decoder[O, A]) assertNextToken(is rune) error {
	if token, err := d.reader.Token(); err != nil {
		return err
	} else if delim, isDelim := token.(Delim); !isDelim {
		return fmt.Errorf("unexpected token %s", string(delim))
	} else if delim != Delim(is) {
		return fmt.Errorf("unexpected token %s", string(delim))
	}

	return nil
}
