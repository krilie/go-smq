package go_smq

import (
	"sync"
	"testing"
)

func TestNewSmq(t *testing.T) {
	mu := sync.RWMutex{}
	mu.RLock()
	mu.RUnlock()
}
