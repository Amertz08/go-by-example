package url

import (
	"sync"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	t.Run("byName", func(t *testing.T) {
		t.Parallel()
		time.Sleep(5 * time.Second)
	})

	t.Run("byInventory", func(t *testing.T) {
		t.Parallel()
		time.Sleep(5 * time.Second)
	})
}

var counter int
var mu sync.Mutex

func getCounter() int {
	mu.Lock()
	defer mu.Unlock()
	return counter
}

func incr() {
	mu.Lock()
	defer mu.Unlock()
	counter++
}

func TestIncr(t *testing.T) {
	t.Parallel()
	t.Run("once", func(t *testing.T) {
		t.Parallel()
		incr()
		if getCounter() != 1 {
			t.Errorf("counter = %d; want 1", counter)
		}
	})

	t.Run("twice", func(t *testing.T) {
		t.Parallel()
		incr()
		incr()
		if getCounter() != 3 {
			t.Errorf("counter = %d; want 3", counter)
		}
	})
}
