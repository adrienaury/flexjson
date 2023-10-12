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
		return nil, fmt.Errorf("%w", err)
	}

	delim, isDelim := token.(Delim)

	switch {
	case isDelim && delim == Delim('{'):
		result, err := d.decodeObject()
		if err != nil {
			return result, err
		}

		return d.closeObject(result)
	case isDelim && delim == Delim('['):
		result, err := d.decodeArray()
		if err != nil {
			return result, err
		}

		return d.closeArray(result)
	case isDelim:
		return nil, fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	}

	return token, nil
}

func (d *Decoder[O, A]) DecodeObject() (O, error) { //nolint:ireturn
	if err := d.assertNextToken('{'); err != nil {
		return *new(O), err
	}

	result, err := d.decodeObject()
	if err != nil {
		return result, err
	}

	return d.closeObject(result)
}

func (d *Decoder[O, A]) decodeObject() (O, error) { //nolint:ireturn
	object, err := d.objMaker()
	if err != nil {
		return object, err
	}

	for d.reader.More() {
		if object, err = d.decodeKeyValue(object); err != nil {
			return object, err
		}
	}

	return object, nil
}

func (d *Decoder[O, A]) decodeKeyValue(obj O) (O, error) { //nolint:ireturn
	key, err := d.reader.Token()
	if err != nil {
		return obj, fmt.Errorf("%w", err)
	}

	keystr, keyIsString := key.(string)
	if !keyIsString {
		return obj, fmt.Errorf("%w %T", ErrInvalidKeyType, key)
	}

	result, err := d.Decode()
	if err != nil {
		return obj, err
	}

	return d.objKeyer(obj, keystr, result)
}

func (d *Decoder[O, A]) DecodeArray() (A, error) { //nolint:ireturn
	if err := d.assertNextToken('['); err != nil {
		return *new(A), err
	}

	array, err := d.decodeArray()
	if err != nil {
		return array, err
	}

	return d.closeArray(array)
}

func (d *Decoder[O, A]) decodeArray() (A, error) { //nolint:ireturn
	array, err := d.arrMaker()
	if err != nil {
		return array, err
	}

	for d.reader.More() {
		item, err := d.Decode()
		if err != nil {
			return array, err
		}

		array, err = d.arrAdder(array, item)
		if err != nil {
			return array, err
		}
	}

	return array, nil
}

func (d *Decoder[O, A]) assertNextToken(is rune) error {
	if token, err := d.reader.Token(); err != nil {
		return fmt.Errorf("%w", err)
	} else if delim, isDelim := token.(Delim); !isDelim {
		return fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	} else if delim != Delim(is) {
		return fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	}

	return nil
}

func (d *Decoder[O, A]) closeObject(obj O) (O, error) { //nolint:ireturn
	if token, err := d.reader.Token(); err != nil {
		return obj, fmt.Errorf("%w", err)
	} else if delim, isDelim := token.(Delim); !isDelim {
		return obj, fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	} else if delim != Delim('}') {
		return obj, fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	}

	return obj, nil
}

func (d *Decoder[O, A]) closeArray(arr A) (A, error) { //nolint:ireturn
	if token, err := d.reader.Token(); err != nil {
		return arr, fmt.Errorf("%w", err)
	} else if delim, isDelim := token.(Delim); !isDelim {
		return arr, fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	} else if delim != Delim(']') {
		return arr, fmt.Errorf("%w %s", ErrUnexpectedToken, string(delim))
	}

	return arr, nil
}
