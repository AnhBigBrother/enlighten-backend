package parser

import (
	"encoding/json"
	"errors"
	"io"
)

func ParseBody(body io.ReadCloser, v interface{}) error {
	if body == nil {
		return errors.New("missing body")
	}
	return json.NewDecoder(body).Decode(v)
}
