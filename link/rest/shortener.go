package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Amertz08/go-by-example/link"
	"github.com/Amertz08/go-by-example/link/kit/hio"
)

// Shortener is a link shortener service that shortens URLs.
type Shortener interface {
	// Shorten shortens a link and returns its key.
	Shorten(ctx context.Context, lnk link.Link) (link.Key, error)
}

// newResponder returns a new HTTP responder with an error handler.
// that maps the errors to the appropriate HTTP status codes.
func newResponder(lg *slog.Logger) hio.Responder {
	err := func(err error) hio.Handler {
		return func(w http.ResponseWriter, r *http.Request) hio.Handler {
			httpError(w, r, lg, err)
			return nil
		}
	}
	return hio.NewResponder(err)
}

// Shorten returns an [http.Handler] that shortens URLs.
func Shorten(lg *slog.Logger, links Shortener) http.Handler {
	with := newResponder(lg)

	return hio.Handler(func(w http.ResponseWriter, r *http.Request) hio.Handler {
		var lnk link.Link
		err := hio.DecodeJSON(hio.MaxBytesReader(w, r.Body, 4_096), &lnk)
		if err != nil {
			return with.Error("decoding: %w: %w", err, link.ErrBadRequest)
		}
		lnk.Key = link.Key(r.PostFormValue("Key"))
		if err := lnk.Validate(); err != nil {
			return with.Error("validating: %w", err)
		}
		key, err := links.Shorten(r.Context(), link.Link{
			Key: link.Key(r.PostFormValue("Key")),
			URL: r.PostFormValue("url"),
		})
		if err != nil {
			return with.Error("shortening: %w", err)
		}

		return with.JSON(http.StatusCreated, map[string]link.Key{"key": key})
	})
}

// Resolver is a link resolver service that resolves shorten link URLs.
type Resolver interface {
	// Resolve resolves a shorten link URL by its key.
	Resolve(ctx context.Context, key link.Key) (link.Link, error)
}

// Resolve returns an [http.Handler] that resolves shorten link URLs.
// It extracts a {kye} from [http.Request] using [http.Request.PathValue].
func Resolve(lg *slog.Logger, links Resolver) http.Handler {
	with := newResponder(lg)

	return hio.Handler(func(w http.ResponseWriter, r *http.Request) hio.Handler {
		lnk, err := links.Resolve(r.Context(), link.Key(r.PathValue("key")))
		if err != nil {
			return with.Error("resolving: %w", err)
		}

		return with.Redirect(http.StatusFound, lnk.URL)
	})
}

func httpError(w http.ResponseWriter, r *http.Request, lg *slog.Logger, err error) {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, link.ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, link.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, link.ErrNotFound):
		code = http.StatusNotFound
	}
	if code == http.StatusInternalServerError {
		lg.ErrorContext(
			r.Context(), "internal", "error", err,
		)
		err = link.ErrInternal
	}
	http.Error(w, err.Error(), code)
}
