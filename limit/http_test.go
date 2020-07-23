package limit

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsRequestLimit(t *testing.T) {

	s := NewMemoryStore()

	now := time.Date(2010, 10, 10, 10, 0, 0, 0, time.UTC)
	timeNow = func() time.Time {
		return now
	}

	// empty test
	limit, tLimit := IsRequestLimit(s, "0", 1, 1)

	assert.Equal(t, false, limit)

	done := make(chan int)

	// add a request per second for ip 0 10X
	for i := 0; i < 10; i++ {
		now = now.Add(time.Second)
		go SaveRequest(s, "0", done)
		<-done
	}

	limit, tLimit = IsRequestLimit(s, "0", 10, 1)

	assert.Equal(t, true, limit)
	assert.Equal(t, now.UnixNano(), tLimit)

	limit, tLimit = IsRequestLimit(s, "0", 11, 1)

	assert.Equal(t, false, limit)

	// add a request per minute for ip 1 10X
	for i := 0; i < 10; i++ {
		now = now.Add(time.Minute)
		go SaveRequest(s, "1", done)
		<-done
	}

	limit, _ = IsRequestLimit(s, "1", 10, 1)

	assert.Equal(t, false, limit)

	limit, _ = IsRequestLimit(s, "1", 11, 1)

	assert.Equal(t, false, limit)

	limit, tLimit = IsRequestLimit(s, "1", 10, 10)

	assert.Equal(t, true, limit)
	assert.Equal(t, now.UnixNano(), tLimit)

	limit, _ = IsRequestLimit(s, "1", 11, 10)

	assert.Equal(t, false, limit)


	// add a request per second for ip 2 100X
	for i := 0; i < 100; i++ {
		now = now.Add(time.Second)
		go SaveRequest(s, "1", done)
		<-done
	}

	lastLimit := now.UnixNano()

	limit, tLimit = IsRequestLimit(s, "1", 100, 60)

	assert.Equal(t, true, limit)
	assert.Equal(t, now.UnixNano(), tLimit)

	// move time 30min ahead
	now = now.Add(time.Minute * 30)

	limit, tLimit = IsRequestLimit(s, "1", 100, 60)

	assert.Equal(t, true, limit)
	assert.Equal(t, lastLimit, tLimit)

	// move time again 30min ahead
	now = now.Add(time.Minute * 30)

	limit, tLimit = IsRequestLimit(s, "1", 100, 60)

	assert.Equal(t, false, limit)

}
