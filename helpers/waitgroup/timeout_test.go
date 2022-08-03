package waitgroup

import (
	"sync"
	"testing"
	"time"
)

func TestWgTimeout(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(5) * time.Second)
	}()

	if !WgTimeout(&wg, time.Duration(1)*time.Second) {
		t.Fatal("WgTimeout didn't timeout after passed Duration")
	}
}
