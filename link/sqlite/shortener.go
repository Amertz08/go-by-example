package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Amertz08/go-by-example/link"
)

// Shortener is a link shortener service that is backed by SQLite.
type Shortener struct {
	db *sql.DB
}

// NewShortener returns a new [Shortener] service.
func NewShortener(db *sql.DB) *Shortener {
	return &Shortener{db: db}
}

// Shorten shortens the URL of the [link.Link] and erturns a [link.Key].
func (s *Shortener) Shorten(
	ctx context.Context, lnk link.Link,
) (link.Key, error) {
	var err error
	if lnk.Key, err = link.Shorten(lnk); err != nil {
		return "", fmt.Errorf("%w: %w", err, link.ErrBadRequest)
	}

	// Persist the link in the database.
	_, err = s.db.ExecContext(ctx, "INSERT INTO links (key, url) VALUES (?, ?)", lnk.Key, lnk.URL)
	if err != nil {
		return "", fmt.Errorf("persisting: %w: %w", err, link.ErrInternal)
	}
	return lnk.Key, nil
}
