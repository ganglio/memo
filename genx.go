package memo

import (
	"sync"
	"time"
)

type MX[T any] func() (T, error)

func (g MX[T]) Memo(r time.Duration) (MX[T], error) {
	data, err := g()
	if err != nil {
		return g, err
	}
	m := struct {
		sync.Mutex
		data            T
		err             error
		lastUpdate      time.Time
		refreshInterval time.Duration
		refreshing      bool
	}{
		data:            data,
		err:             err,
		lastUpdate:      time.Now(),
		refreshInterval: r,
		refreshing:      false,
	}
	return func() (T, error) {
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
						m.err = nil
						m.lastUpdate = time.Now()
					} else {
						m.err = err
					}
					m.refreshing = false
					m.Unlock()
				}()
			}
		}
		return m.data, m.err
	}, nil
}
