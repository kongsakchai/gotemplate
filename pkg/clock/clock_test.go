package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func TestNow(t *testing.T) {
	c := New()
	now := c.Now()
	assert.False(t, now.IsZero())

	expected := time.Now()
	diff := expected.Sub(now)
	assert.Less(t, diff, time.Second)
}
