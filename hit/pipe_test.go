package hit

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"
)

func TestDispatchErrorThresholdEndsTest(t *testing.T) {
	errThreshold := 1
	inputChannel := produce(10, &http.Request{})
	send := func(*http.Request) Result {
		var err error
		if rand.IntN(2) == 1 {
			err = fmt.Errorf("test error")
		}
		return Result{
			Status:   200,
			Bytes:    0,
			Duration: 0,
			Error:    err,
		}
	}
	results := dispatch(inputChannel, 1, errThreshold, send)

	errCount := 0
	resultCount := 0
	for result := range results {
		resultCount++
		if result.Error != nil {
			errCount++
		}
	}

	if resultCount == 10 {
		t.Fatalf("expected at least 1 error, got none")
	}

	if errCount > errThreshold {
		t.Errorf("expected only %d error, got %d", errThreshold, errCount)
	}
}

func TestDispatchErrorThresholdZeroUsesAllRequests(t *testing.T) {
	errThreshold := 0
	requestCount := 10
	inputChannel := produce(requestCount, &http.Request{})
	send := func(*http.Request) Result {
		var err error
		if rand.IntN(2) == 1 {
			err = fmt.Errorf("test error")
		}
		return Result{
			Status:   200,
			Bytes:    0,
			Duration: 0,
			Error:    err,
		}
	}
	results := dispatch(inputChannel, 1, errThreshold, send)

	// TODO: Is there a better way to check this?
	resultCount := 0
	for range results {
		resultCount++
	}

	if resultCount != requestCount {
		t.Errorf("expected 10 results, got %d", resultCount)
	}
}
