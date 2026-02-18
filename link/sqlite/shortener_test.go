package sqlite

import (
	"errors"
	"testing"

	"github.com/Amertz08/go-by-example/link"
)

func TestShortenerShorten(t *testing.T) {
	t.Parallel()

	lnk := link.Link{
		Key: "foo",
		URL: "https://new.link",
	}

	shortener := NewShortener(DialTestDB(t))

	// Shortens a link.
	key, err := shortener.Shorten(t.Context(), lnk)
	if err != nil {
		t.Fatalf("shortening link: %v", err)
	}
	if key != lnk.Key {
		t.Errorf("expected key %q, got %q", lnk.Key, key)
	}

	// Disallows shortening a link with a duplicate key.
	_, err = shortener.Shorten(t.Context(), lnk)
	if !errors.Is(err, link.ErrConflict) {
		t.Fatalf("expected error %w, got %v", link.ErrConflict, err)
	}
}
