package memo

import (
	"sync"
	"time"
)

// Memo creates a new cached variable with a given refresh interval
func Memo(g func() any, r time.Duration) func() any {
	m := struct {
		sync.Mutex
		data            interface{}
		lastUpdate      time.Time
		refreshInterval time.Duration
		refreshing      bool
	}{
		data:            g(),
		lastUpdate:      time.Now(),
		refreshInterval: r,
		refreshing:      false,
	}
	return func() any {
		m.Lock()
		defer m.Unlock()
		if time.Since(m.lastUpdate) > m.refreshInterval {
			if !m.refreshing {
				m.refreshing = true
				go func() {
					data := g()
					m.Lock()
					m.data = data
					m.lastUpdate = time.Now()
					m.refreshing = false
					m.Unlock()
				}()
			}
		}
		return m.data
	}
}

func MemoX(g func() (any, error), r time.Duration) func() any {
	data, err := g()
	if err != nil {
		return func() any {
			return err
		}
	}
	m := struct {
		sync.Mutex
		data            interface{}
		lastUpdate      time.Time
		refreshInterval time.Duration
		refreshing      bool
	}{
		data:            data,
		lastUpdate:      time.Now(),
		refreshInterval: r,
		refreshing:      false,
	}
	return func() interface{} {
		m.Lock()
		defer m.Unlock()
		if time.Since(m.lastUpdate) > m.refreshInterval {
			if !m.refreshing {
				m.refreshing = true
				go func() {
					data, err := g()
					m.Lock()
					if err == nil {
						m.data = data
					}
					m.lastUpdate = time.Now()
					m.refreshing = false
					m.Unlock()
				}()
			}
		}
		return m.data
	}
}
