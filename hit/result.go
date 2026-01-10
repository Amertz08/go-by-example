package hit

import (
	"iter"
	"time"
)

// Result is performance metrics of a single [http.Request]
type Result struct {
	Status   int
	Bytes    int64
	Duration time.Duration
	Error    error
}

// Results is an iterator for [Result] values
type Results iter.Seq[Result]
