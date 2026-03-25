package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimerNew(t *testing.T) {
	t.Run("should create a new timer with the specified location", func(t *testing.T) {
		// arrange
		loc, err := time.LoadLocation("UTC")
		require.NoError(t, err)

		// act
		tm := New(loc)
		t.Log(tm.Now())

		// assert
		assert.Equal(t, loc, tm.loc)
		assert.Equal(t, tm.Now().Location().String(), loc.String())
	})

	t.Run("should create a new timer with local location when loc is nil", func(t *testing.T) {
		// act
		tm := New(nil)
		t.Log(tm.Now())

		// assert
		assert.Equal(t, time.Local, tm.loc)
		assert.Equal(t, tm.Now().Location().String(), time.Local.String())
	})
}
