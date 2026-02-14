package hit

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

var commonStatusCodes = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusInternalServerError,
}

var errorMap = map[int]error{
	http.StatusBadRequest:          errors.New("bad request"),
	http.StatusInternalServerError: errors.New("internal server error"),
}

func RandomStatusCode() int {
	return commonStatusCodes[rand.IntN(len(commonStatusCodes))]
}

// Send sends an HTTP request and returns a performance [Result]
func Send(_ *http.Client, _ *http.Request) Result {
	const roundTripTime = 100 * time.Millisecond

	time.Sleep(roundTripTime)

	statusCode := RandomStatusCode()

	var err error
	if statusCode != http.StatusOK {
		err = errorMap[statusCode]
	}

	return Result{
		Status:   statusCode,
		Bytes:    10,
		Duration: roundTripTime,
		Error:    err,
	}
}

// SendN sends N requests using [Send].
// It returns a single-use [Results] iterator that
// pushes a [Result] for each [http.Request] sent.
func SendN(n int, req *http.Request, opts Options) (Results, error) {
	opts = withDefaults(opts)
	if n <= 0 {
		return nil, fmt.Errorf("n musts be positive: got %d", n)
	}

	results := runPipeline(n, req, opts)

	return func(yield func(Result) bool) {
		for result := range results {
			if !yield(result) {
				return
			}
		}
	}, nil
}
