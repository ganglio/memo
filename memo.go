package memo

import (
	"sync"
	"time"
)

// Memo creates a new cached variable with a given refresh interval
func Memo(g func() interface{}, r time.Duration) func() interface{} {
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
	return func() interface{} {
		m.Lock()
		defer m.Unlock()
		if time.Now().Sub(m.lastUpdate) > m.refreshInterval {
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

func MemoX(g func() (interface{}, error), r time.Duration) func() interface{} {
	data, err := g()
	if err != nil {
		return func() interface{} {
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
		if time.Now().Sub(m.lastUpdate) > m.refreshInterval {
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
