package memo

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemo(t *testing.T) {

	f := Memo(func() interface{} {
		return fmt.Sprintf("Pesco %s", time.Now().Format("15:04:05"))
	}, time.Second)

	now := fmt.Sprintf("Pesco %s", time.Now().Format("15:04:05"))

	assert.Equal(t, now, f())
	time.Sleep(3 * time.Second)
	assert.Equal(t, now, f()) // cached value and trigger refresh
	time.Sleep(250 * time.Millisecond)
	assert.NotEqual(t, now, f()) // refreshed value
}

func TestMemoX(t *testing.T) {
	v := 0
	f := MemoX(func() (interface{}, error) {
		v = v + 1
		if v >= 3 && v <= 7 {
			return v, errors.New("Raising exception!")
		}
		return v, nil
	}, 10*time.Millisecond)

	p := -1
	for i := 0; i < 20; i++ {
		p = f().(int)
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, v, 20)
	assert.Equal(t, p, 19)
}

func TestStampede(t *testing.T) {
	v := 0
	f := Memo(func() interface{} {
		v = v + 1
		return v
	}, time.Second)

	t0 := time.Now()

	c := 0
	o := 0
	for time.Now().Sub(t0) < 10*time.Second {
		c = c + 1
		o = f().(int)
	}

	assert.Less(t, o, c)
	assert.Equal(t, 10, o)
}

func BenchmarkMemo(b *testing.B) {
	v := 0
	f := Memo(func() interface{} {
		v = v + 1
		return v
	}, time.Second)

	for n := 0; n < b.N; n++ {
		f()
	}
}
