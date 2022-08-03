package befehl

import (
	"sync/atomic"
)

type queue struct {
	count int64
}

func (q *queue) New(hostCnt int64) *queue {
	return &queue{
		count: hostCnt,
	}
}

func (q *queue) decrementCounter(total int) int64 {
	return atomic.AddInt64(&q.count, -1)
}
