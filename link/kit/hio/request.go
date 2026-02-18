package hio

import (
	"encoding/json"
	"fmt"
	"io"
)

// DecodeJSON reads and decodes JSON.
func DecodeJSON(from io.Reader, to any) error {
	data, err := io.ReadAll(from)
	if err != nil {
		return fmt.Errorf("reading request body: %w", err)
	}
	if err := json.Unmarshal(data, to); err != nil {
		return fmt.Errorf("decoding json: %w", err)
	}
	// If the target implements the Validate interface, validate it.
	v, ok := to.(interface{ Validate() error })
	if ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validating: %w", err)
		}
	}
	return nil
}
