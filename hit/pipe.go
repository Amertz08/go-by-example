package hit

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func produce(
	ctx context.Context,
	n int,
	req *http.Request,
) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)
		for range n {
			select {
			case <-ctx.Done():
				return
			case out <- req.Clone(ctx):
			}
		}
	}()

	return out
}

func throttle(
	ctx context.Context,
	in <-chan *http.Request,
	delay time.Duration,
) <-chan *http.Request {
	out := make(chan *http.Request)

	go func() {
		defer close(out)
		t := time.NewTicker(delay)
		for r := range in {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				out <- r
			}
		}
	}()
	return out
}

func dispatch(
	ctx context.Context,
	in <-chan *http.Request,
	concurrency int,
	errorThreshold int,
	send SendFunc,
) <-chan Result {
	out := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	errCount := 0
	for range concurrency {
		go func() {
			defer wg.Done()
			for req := range in {
				select {
				case <-ctx.Done():
					return
				default:
					if errorThreshold != 0 && errCount >= errorThreshold {
						return
					}
					resp := send(req)
					if resp.Error != nil {
						errCount++
					}
					out <- resp
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func runPipeline(ctx context.Context, n int, req *http.Request, opts Options) <-chan Result {
	requests := produce(ctx, n, req)
	if opts.RPS > 0 {
		requests = throttle(ctx, requests, time.Second/time.Duration(opts.RPS))
	}
	return dispatch(ctx, requests, opts.Concurrency, opts.ErrorThreshold, opts.Send)
}
