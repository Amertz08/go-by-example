package link

import (
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
func (lnk Link) Validate() error {
	if err := lnk.Key.Validate(); err != nil {
		return fmt.Errorf("key: %w", err)
	}
	u, err := url.ParseRequestURI(lnk.URL)
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
