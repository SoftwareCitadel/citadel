package util

import (
	"bytes"
	"encoding/json"
)

// EncodeJSON encodes a struct to JSON.
func EncodeJSON(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecodeJSON decodes JSON to a struct.
func DecodeJSON(data []byte, v any) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}
