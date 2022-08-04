package queue

import (
	"sync/atomic"
)

type Queue struct {
	count int64
}

func (q *Queue) New(hostCnt int64) *Queue {
	return &Queue{
		count: hostCnt,
	}
}

func (q *Queue) DecrementCounter() int64 {
	return atomic.AddInt64(&q.count, -1)
}
