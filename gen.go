package memo

import (
	"sync"
	"time"
)

type M[T any] func() T

func (g M[T]) Memo(r time.Duration) M[T] {
	m := struct {
		sync.Mutex
		data            T
		lastUpdate      time.Time
		refreshInterval time.Duration
		refreshing      bool
	}{
		data:            g(),
		lastUpdate:      time.Now(),
		refreshInterval: r,
		refreshing:      false,
	}
	return func() T {
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
