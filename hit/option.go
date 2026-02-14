package hit

import (
	"net/http"
	"time"
)

// SendFunc is a type of function that sends an
// [http.Request] and returns a [Result].
type SendFunc func(*http.Request) Result

// Options defines the options for sending requests.
// Uses default options for unset options.
type Options struct {
	// concurrency is the number of concurrent requests to send.
	// Default: 1
	Concurrency int

	// RPS is the requests to send per second.
	// Default: 0 (no rate limiting)
	RPS int

	// ErrorThreshold is the number of errors to tolerate before exiting.
	// Default: 0 (never exit)
	ErrorThreshold int

	// Send processes requests.
	// Default Uses [Send].
	Send SendFunc

	// Client is the HTTP client to use.
	// Default Uses [http.DefaultClient].
	Client *http.Client
}

func Defaults() Options {
	return withDefaults(Options{})
}

func withDefaults(o Options) Options {
	if o.Concurrency == 0 {
		o.Concurrency = 1
	}
	if o.Send == nil {
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: o.Concurrency,
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
			Timeout:       30 * time.Second,
		}
		o.Send = func(r *http.Request) Result { return Send(client, r) }
	}
	if o.Client == nil {
		o.Client = http.DefaultClient
	}
	return o
}
