package counter

import (
	"expvar"
	"strconv"
	"sync"
)

var c *Counter

type Counter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *Counter) Get() int {
	return c.cnt
}

func (c *Counter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *Counter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

func New(name string) *Counter {
	c = &Counter{m: &sync.RWMutex{}}
	expvar.Publish(name, c)
	return c
}
