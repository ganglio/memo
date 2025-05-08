package memo

import (
	"sync"
	"time"
)

type Gen[T any] func() T
type GenX[T any] func() (T, error)

func (g Gen[T]) Memo(r time.Duration) Gen[T] {
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

func (g GenX[T]) Memo(r time.Duration) (Gen[T], error) {
	data, err := g()
	if err != nil {
		return nil, err
	}
	m := struct {
		sync.Mutex
		data            T
		lastUpdate      time.Time
		refreshInterval time.Duration
		refreshing      bool
	}{
		data:            data,
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
					data, err := g()
					m.Lock()
					if err == nil {
						m.data = data
						m.lastUpdate = time.Now()
					}
					m.refreshing = false
					m.Unlock()
				}()
			}
		}
		return m.data
	}, nil
}
