package cascade

import "sync/atomic"

type IOCounter struct {
	count int64
}

func NewIOCounter() *IOCounter { return &IOCounter{} }

func (c *IOCounter) Increment() { atomic.AddInt64(&c.count, 1) }
func (c *IOCounter) Reset()     { atomic.StoreInt64(&c.count, 0) }
func (c *IOCounter) Count() int64 { return atomic.LoadInt64(&c.count) }
