package hit

import (
	"iter"
	"net/http"
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

// Summary is the summary of [Result] values.
type Summary struct {
	Requests int
	Errors   int
	Bytes    int64
	RPS      float64
	Duration time.Duration
	Fastest  time.Duration
	Slowest  time.Duration
	Success  float64
}

// Summarize returns a [Summary] of the [Results].
func Summarize(results Results) Summary {
	var s Summary
	if results == nil {
		return s
	}

	started := time.Now()
	for r := range results {
		s.Requests++
		s.Bytes += r.Bytes

		if r.Error != nil || r.Status != http.StatusOK {
			s.Errors++
		}
		if s.Fastest == 0 {
			s.Fastest = r.Duration
		}
		if r.Duration < s.Fastest {
			s.Fastest = r.Duration
		}
		if r.Duration > s.Slowest {
			s.Slowest = r.Duration
		}
		if r.Duration > s.Slowest {
			s.Slowest = r.Duration
		}
	}
	if s.Requests > 0 {
		s.Success = (float64(s.Requests-s.Errors) / float64(s.Requests)) * 100
	}
	s.Duration = time.Since(started)
	s.RPS = float64(s.Requests) / s.Duration.Seconds()

	return s
}
