package memo

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenMemo(t *testing.T) {
	f := M[string](func() string {
		return fmt.Sprintf("Pesco %s", time.Now().Format("15:04:05"))
	}).Memo(time.Second)

	now := fmt.Sprintf("Pesco %s", time.Now().Format("15:04:05"))

	assert.Equal(t, now, f())
	time.Sleep(3 * time.Second)
	assert.Equal(t, now, f()) // cached value and trigger refresh
	time.Sleep(250 * time.Millisecond)
	assert.NotEqual(t, now, f()) // refreshed value
}

func TestGenMemoX(t *testing.T) {
	v := 0
	f, err := MX[int](func() (int, error) {
		v = v + 1
		if v >= 3 && v <= 7 {
			return v, errors.New("Raising exception!")
		}
		return v, nil
	}).Memo(10 * time.Millisecond)

	assert.NoError(t, err)

	p := -1
	for i := 0; i < 20; i++ {
		p = f()
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, v, 20)
	assert.Equal(t, p, 19)
}

func TestGenMemoX2(t *testing.T) {
	_, err := MX[int](func() (int, error) {
		return -1, errors.New("Raising exception!")
	}).Memo(10 * time.Millisecond)

	assert.Error(t, err)
}

func TestGenStampede(t *testing.T) {
	v := 0
	f := M[int](func() int {
		v = v + 1
		return v
	}).Memo(time.Second)

	t0 := time.Now()

	c := 0
	o := 0
	for time.Since(t0) < 10*time.Second {
		c = c + 1
		o = f()
	}

	assert.Less(t, o, c)
	assert.Equal(t, 10, o)
}

func BenchmarkGenMemo(b *testing.B) {
	v := 0
	f := M[int](func() int {
		v = v + 1
		return v
	}).Memo(time.Second)

	for n := 0; n < b.N; n++ {
		f()
	}
}
