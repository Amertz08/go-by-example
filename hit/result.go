package hit

import "time"

// Result is performance metrics of a single [http.Request]
type Result struct {
	Status   int
	Bytes    int64
	Duration time.Duration
	Error    error
}
