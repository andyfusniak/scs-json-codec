package codec

import (
	"bytes"
	"encoding/json"
	"time"
)

// JSONCodec is used for encoding/decoding session data to and from a byte
// slice using the encoding/json package.
type JSONCodec struct{}

// Encode converts a session deadline and values into a byte slice.
func (JSONCodec) Encode(deadline time.Time, values map[string]interface{}) ([]byte, error) {
	aux := &struct {
		Deadline time.Time
		Values   map[string]interface{}
	}{
		Deadline: deadline,
		Values:   values,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&aux); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Decode converts a byte slice into a session deadline and values.
func (JSONCodec) Decode(b []byte) (time.Time, map[string]interface{}, error) {
	aux := &struct {
		Deadline time.Time
		Values   map[string]interface{}
	}{}

	r := bytes.NewReader(b)
	if err := json.NewDecoder(r).Decode(&aux); err != nil {
		return time.Time{}, nil, err
	}

	return aux.Deadline, aux.Values, nil
}
