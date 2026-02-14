package hit

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"
)

func TestDispatch(t *testing.T) {
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
	for result := range results {
		if result.Error != nil {
			errCount++
		}
	}

	if errCount > errThreshold {
		t.Errorf("expected only %d error, got %d", errThreshold, errCount)
	}
}
