package v1

import (
	"encoding/json"

	"github.com/7sDream/geko"
)

type Token = json.Token

type Delim = json.Delim

type TokenReader interface {
	More() bool
	Token() (Token, error)
}

type Object = geko.Object
