package req

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseBody(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}
