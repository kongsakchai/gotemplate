package mockutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	t.Run("should create a new timer", func(t *testing.T) {
		// act
		timer := NewTimer(t)

		// assert
		assert.NotNil(t, timer)
	})

	t.Run("should return the expected time when Now is called", func(t *testing.T) {
		// arrange
		timer := NewTimer(t)
		expectedTime := time.Date(2024, time.June, 1, 12, 0, 0, 0, time.UTC)
		timer.On("Now").Return(expectedTime)

		// act
		actualTime := timer.Now()

		// assert
		assert.Equal(t, expectedTime, actualTime)
		timer.AssertExpectations(t)
	})
}
