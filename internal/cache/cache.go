package cache

import "strings"

type Cache struct {
	m map[string]any
}

func New() *Cache {
	return &Cache{
		m: map[string]any{},
	}
}

func (c *Cache) Set(cmd string, resp any) {
	if c == nil {
		return
	}
	c.m[cmd] = resp
}

func (c *Cache) Get(cmd string) any {
	if c == nil {
		return nil
	}

	return c.m[cmd]
}

func (c *Cache) Update(namespace string) {
	if c == nil {
		return
	}
	for k := range c.m {
		if strings.HasPrefix(k, namespace) {
			delete(c.m, k)
		}
	}
}
