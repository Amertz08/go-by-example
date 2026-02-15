package link

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Link struct {
	// URL is the original URL of the link.
	URL string
	// Key is the shortened key of the URL.
	Key Key
}

// Key is the shortened key of a [Link] URL.
type Key string

// String returns the key without leading or trailing spaces.
func (k Key) String() string { return strings.TrimSpace(string(k)) }

// Empty reports whether the [Key] is empty.
func (k Key) Empty() bool { return k.String() == "" }

// Validate validates the [Link].
func (l Link) Validate() error {
	if err := l.Key.Validate(); err != nil {
		return fmt.Errorf("key: %w", err)
	}
	u, err := url.ParseRequestURI(l.URL)
	if err != nil {
		return err
	}
	if u.Host == "" {
		return errors.New("empty host")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	return nil
}

// Validate validates the [Key].
func (k Key) Validate() error {
	// We use generated (8-hex-character) by default, but allow
	// user-defined keys up to 16 characters for convenience.
	const maxKeyLen = 16
	if len(k.String()) > maxKeyLen {
		return fmt.Errorf("too long (max %d)", maxKeyLen)
	}
	return nil
}

// Shorten shortens the [Link] URL and generates a new [Key]
// if the [Key] is empty. Otherwise, it returns the same [Key].
// It returns an error if the [Link] is invalid.
func Shorten(l Link) (Key, error) {
	if l.Key.Empty() {
		sum := sha256.Sum256([]byte(l.URL))
		l.Key = Key(base64.RawURLEncoding.EncodeToString(sum[:6]))
	}
	if err := l.Validate(); err != nil {
		return "", fmt.Errorf("validating: %w")
	}
	return l.Key, nil
}
