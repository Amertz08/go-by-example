package hit

import "net/http"

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
}

func Defaults() Options {
	return withDefaults(Options{})
}

func withDefaults(o Options) Options {
	if o.Concurrency == 0 {
		o.Concurrency = 1
	}
	if o.Send == nil {
		o.Send = func(r *http.Request) Result {
			return Send(http.DefaultClient, r)
		}
	}
	return o
}
