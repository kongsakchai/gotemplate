package null

import (
	"bytes"
	"encoding/json"
)

type N[T any] struct {
	V     T
	Valid bool
}

// MarshalJSON implements json.Marshaler for N[T].
// If Valid is false, it marshals as JSON null; otherwise it marshals V.
func (n N[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.V)
}

// UnmarshalJSON implements json.Unmarshaler for N[T].
// JSON null sets Valid to false; any other value is unmarshaled into V and sets Valid to true.
func (n *N[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		var zero T
		n.V = zero
		n.Valid = false
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.V = v
	n.Valid = true
	return nil
}
