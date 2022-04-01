package internal

import (
	"sync"
	"time"

	"github.com/prometheus/common/log"
)

type Refresher[T any] interface {
	Interval() time.Duration
	Refresh() ([]T, error)
}

type Cache[T any] struct {
	Logger    log.Logger
	Refresher Refresher[T]

	mutex     sync.RWMutex
	timestamp time.Time
	data      []T
}

func (c *Cache[T]) IsValid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return time.Now().Sub(c.timestamp) < c.Refresher.Interval()
}

func (c *Cache[T]) Timestamp() time.Time {
	if !c.IsValid() {
		go c.refresh()
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.timestamp
}

func (c *Cache[T]) Data() []T {
	if !c.IsValid() {
		go c.refresh()
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.data
}

func (c *Cache[T]) refresh() {
	d, err := c.Refresher.Refresh()
	if err != nil {
		c.Logger.Debugf("cache: failed to refresh: %s", err)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.timestamp = time.Now()
	c.data = d
}

func NewCache[T any](l log.Logger, r Refresher[T]) *Cache[T] {
	return &Cache[T]{
		Logger:    l,
		Refresher: r,
	}
}
