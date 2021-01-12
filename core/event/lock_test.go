package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventLock(t *testing.T) {
	e := &Event{
		EventID:     "123",
		PublishTime: time.Now(),
		Data:        "TEST_TEST_1234",
	}

	err := e.Lock(2 * time.Minute)
	assert.Nil(t, err, "should not throw an error")
	assert.True(t, e.IsLocked, "should be locked")

	err = e.Lock(2 * time.Minute)
	assert.Error(t, err, "should throw an error, could not be locked")
	assert.True(t, e.IsLocked, "should be already locked")

	e.Unlock()
	assert.False(t, e.IsLocked, "should be already locked")
}

func TestEventLockExpire(t *testing.T) {
	e := &Event{
		EventID:     "123",
		PublishTime: time.Now(),
		Data:        "TEST_TEST_1234",
	}
	ttl := 200 * time.Microsecond

	err := e.Lock(ttl)
	assert.Nil(t, err, "should not throw an error")
	assert.True(t, e.IsLocked, "should be locked")

	time.Sleep(ttl)
	assert.False(t, e.IsLocked, "should not be locked anymore")
}

func TestEventLockExpireOn0TTL(t *testing.T) {
	e := &Event{
		EventID:     "123",
		PublishTime: time.Now(),
		Data:        "TEST_TEST_1234",
	}

	err := e.Lock(0)
	assert.Nil(t, err, "should not throw an error")

	time.Sleep(100 * time.Microsecond)
	// The lock should just be gone after the default TTL of 1 minute
	assert.True(t, e.IsLocked, "should be locked, not directly unlocked")
}

func TestEventUnlockOnUnlocked(t *testing.T) {
	e := &Event{
		EventID:     "123",
		PublishTime: time.Now(),
		Data:        "TEST_TEST_1234",
	}

	e.Unlock()
	assert.False(t, e.IsLocked, "should be still unlocked")
}
