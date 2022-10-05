package clientcache

import (
	"net/http"
	"sync"
	"time"
)

type ClientCache interface {
	Get(interface{}) *http.Client
	Put(interface{}, *http.Client)
}

type clientCache struct {
	mu  sync.Mutex
	now func() time.Time
	ttl time.Duration
	m   map[interface{}]*v
}

type v struct {
	expires time.Time
	cli     *http.Client
}

// NewClientCache returns a new ClientCache
func NewClientCache(ttl time.Duration) ClientCache {
	return &clientCache{
		now: time.Now,
		ttl: ttl,
		m:   map[interface{}]*v{},
	}
}

// call holding c.mu
func (c *clientCache) expire() {
	now := c.now()
	for k, v := range c.m {
		if now.After(v.expires) {
			v.cli.CloseIdleConnections()
			delete(c.m, k)
		}
	}
}

func (c *clientCache) Get(k interface{}) (cli *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v := c.m[k]; v != nil {
		v.expires = c.now().Add(c.ttl)
		cli = v.cli
	}

	c.expire()

	return
}

func (c *clientCache) Put(k interface{}, cli *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[k] = &v{
		expires: c.now().Add(c.ttl),
		cli:     cli,
	}

	c.expire()
}
