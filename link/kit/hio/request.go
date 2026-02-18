package hio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// MaxBytesReader is like [http.MaxBytesReader].
func MaxBytesReader(w http.ResponseWriter, rc io.ReadCloser, max int64) io.ReadCloser {
	type unwrapper interface {
		Unwrap() http.ResponseWriter
	}
	for {
		v, ok := w.(unwrapper)
		if !ok {
			break
		}
		w = v.Unwrap()
	}

	return http.MaxBytesReader(w, rc, max)
}
