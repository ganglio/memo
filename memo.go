package memo

import (
	"sync"
	"time"
)

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
