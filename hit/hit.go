package hit

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Send sends an HTTP request and returns a performance [Result]
func Send(client *http.Client, req *http.Request) Result {
	started := time.Now()
	var (
		bytes int64
		code  int
	)

	resp, err := client.Do(req)
	if err == nil {
		{
			defer resp.Body.Close()
			code = resp.StatusCode
			bytes, err = io.Copy(ioutil.Discard, resp.Body)
		}
	}
	return Result{
		Duration: time.Since(started),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
}

// SendN sends N requests using [Send].
// It returns a single-use [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(ctx context.Context, n int, req *http.Request, opts Options) (Results, error) {
	opts = withDefaults(opts)
	if n <= 0 {
		return nil, fmt.Errorf("n musts be positive: got %d", n)
	}

	ctx, cancel := context.WithCancel(ctx)

	results := runPipeline(ctx, n, req, opts)

	return func(yield func(Result) bool) {
		defer cancel()
		for result := range results {
			if !yield(result) {
				return
			}
		}
	}, nil
}
